package user

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/scofieldpeng/go-libs/tools"
	"github.com/scofieldpeng/todo/api/libs/auth"
	"github.com/scofieldpeng/todo/api/libs/common"
	"github.com/scofieldpeng/todo/api/libs/email"
	"github.com/scofieldpeng/todo/api/libs/redis/find"
	"github.com/scofieldpeng/todo/api/models/user"
	"log"
	"net/http"
	"regexp"
	"time"
)

// loginData 登录请求数据
type loginData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// registerData 注册请求数据
type registerData struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// findData 找回密码请求数据
type findData struct {
	Type     string `json:"type"`
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

// resetData 充值密码请求数据
type resetData struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

// isValidUserName 检查是否是有效的用户名
func validUserName(username string) error {
	if matched, err := regexp.Match("^[a-zA-z_-]{6,20}$", []byte(username)); err != nil || !matched {
		return errors.New("非法用户名")
	}

	return nil
}

// checkPassword 检查密码是否正确
func validPassword(password string) error {
	if matched, err := regexp.Match("^[a-zA-z0-9]{1,32}$", []byte(password)); err != nil || !matched {
		return errors.New("非法密码")
	}
	return nil
}

// Login 用户登录
func Login(ctx echo.Context) error {
	var logindata loginData
	if err := common.GetBodyStruct(ctx, &logindata); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 201, "请输入正确的数据")
	}

	if err := validUserName(logindata.UserName); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 202, err.Error())
	}

	if err := validPassword(logindata.Password); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 203, err.Error())
	}

	// 检查用户名是否存在
	userModel := user.New()
	userModel.UserName = logindata.UserName
	if exsit, err := userModel.Get(); err != nil {
		log.Println("获取用户数据失败!用户名:", logindata.UserName, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 204)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 205, "用户不存在")
	}

	// 检查用户密码
	if userModel.Password != tools.Md5([]byte(logindata.Password+userModel.Salt)) {
		return common.BackError(ctx, http.StatusBadRequest, 206, "用户密码不正确")
	}

	if _, err := auth.SaveToken(ctx, userModel.UserID); err != nil {
		log.Println("生成login token时失败!用户名:", userModel.UserName, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 207)
	}

	lastLogin := userModel.LastLogin

	// 更新用户上次登录时间
	userModel.LastLogin = int(time.Now().Unix())
	if _, err := userModel.Update(); err != nil {
		log.Println("更新用户最近一次登录时间时出错,错误原因:", err.Error())
	}
	userModel.LastLogin = lastLogin

	return ctx.JSON(http.StatusOK, userModel)
}

// Register 用户注册
// 用户注册接口提交邮箱,用户名,密码即可,前期暂时不使用验证码进行验证
func Register(ctx echo.Context) error {
	var data registerData
	if err := common.GetBodyStruct(ctx, &data); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 201, "上传的数据格式不正确")
	}

	// 检查用户名和邮箱是否符合要求
	if err := validUserName(data.UserName); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 202, err.Error())
	}
	if !tools.IsEmail(data.Email) {
		return common.BackError(ctx, http.StatusBadRequest, 203, "请输入正确的邮箱")
	}
	if err := validPassword(data.Password); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 204, err.Error())
	}

	// 检查用户名和邮箱是否已经被注册
	userModel := user.New()
	userModel.UserName = data.UserName
	if exsit, err := userModel.Get(); err != nil {
		log.Println("获取用户信息时失败,获取的用户user_name:", data.UserName, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 205)
	} else if exsit {
		return common.BackError(ctx, http.StatusBadRequest, 206, "该用户名已经注册")
	}

	userModel = user.New()
	userModel.Email = data.Email
	if exsit, err := userModel.Get(); err != nil {
		log.Println("获取用户信息时失败,获取的用户邮箱:", data.Email, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 207)
	} else if exsit {
		return common.BackError(ctx, http.StatusBadRequest, 208, "该邮箱已经注册")
	}

	// 生成密码
	userModel.Salt = tools.RandomString()
	userModel.Password = tools.Md5([]byte(data.Password + userModel.Salt))
	userModel.CreateTime = int(time.Now().Unix())
	userModel.LastLogin = userModel.CreateTime

	userModel.UserName = data.UserName
	userModel.Email = data.Email
	if _, err := userModel.Insert(); err != nil {
		log.Printf("创建用户失败!创建用户数据,%#v,错误原因:%s\n", userModel, err.Error())
		return common.BackServerError(ctx, 209)
	}

	return ctx.JSON(http.StatusOK, userModel)
}

// Find 用户找回密码
// 用户找回密码填写找回密码的请求方式
func Find(ctx echo.Context) error {
	var data findData
	if err := common.GetBodyStruct(ctx, &data); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 201, "请上传正确的数据")
	}

	userModel := user.New()
	switch data.Type {
	case "email":
		if !tools.IsEmail(data.Email) {
			return common.BackError(ctx, http.StatusBadRequest, 202, "请上传正确的邮箱")
		}
		userModel.Email = data.Email
		break
	case "user_name":
		if err := validUserName(data.UserName); err != nil {
			return common.BackError(ctx, http.StatusBadRequest, 203, err.Error())
		}
		userModel.UserName = data.UserName
		break
	default:
		return common.BackError(ctx, http.StatusBadRequest, 204, "请上传正确的数据")
	}

	if exist, err := userModel.Get(); err != nil {
		log.Printf("获取用户信息失败!要查找的用户信息:%v,错误原因:%s\n", userModel, err.Error())
		return common.BackServerError(ctx, 205)
	} else if !exist {
		return common.BackError(ctx, http.StatusBadRequest, 206, "该账号没有注册")
	}

	// 获取该用户上次请求时间
	lastTime, err := find.GetLastTime(userModel.UserID)
	if err != nil {
		log.Println("获取用户最近一次请求find接口时间出错,用户id:", userModel.UserID, ",出错原因:", err.Error())
		return common.BackServerError(ctx, 207)
	}

	if lastTime != 0 && lastTime < time.Now().Unix() {
		return common.BackError(ctx, http.StatusBadRequest, 208, "请求时间间隔过短,请稍后再试!")
	}

	// 生成token,时间戳(纳秒）+用户user_name+email+用户当前salt，拼接后sha1加密
	token := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%d-%s-%s-%s", time.Now().Nanosecond(), userModel.UserName, userModel.Email, userModel.Salt))))

	// 将token写入到redis中,并且记录上次邮件找回密码时间戳
	if err := find.SetToken(userModel.UserID, token); err != nil {
		log.Println("写入find_token失败,用户id:", userModel.UserID, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 209)
	}

	// 发送邮件
	emailIns := email.New()
	if err := emailIns.SetTo(userModel.Email).SetTpl(email.Find_Pwd_Tpl).SetTplVals(map[string]string{
		"token": token,
	}).Send(); err != nil {
		log.Printf("发送邮件失败!发送数据:%#v,失败原因:%s\n", emailIns, err.Error())
		return common.BackServerError(ctx, 210)
	}

	return common.BackOk(ctx)
}

// ResetPwd 用户重置密码接口
// 通过判断是否存在重置token,然后根据token进行重置密码操作
func ResetPwd(ctx echo.Context) error {
	var data resetData
	if err := common.GetBodyStruct(ctx, &data); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 201, "服务器内部错误")
	}
	if data.Token == "" {
		return common.BackError(ctx, http.StatusBadRequest, 202, "请传入token值")
	}
	if err := validPassword(data.Password); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 203, "请传入密码")
	}

	// 检查token是否有效
	userid, err := find.GetUserIDFromToken(data.Token)
	if err != nil {
		log.Println("获取redis值中的token失败,token:", data.Token, ",失败原因:", err.Error())
		return common.BackServerError(ctx, 204)
	}
	if userid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 205, "没有找到该用户")
	}

	// 查询该用户是否存在
	userModel := user.New()
	userModel.UserID = userid
	if exsit, err := userModel.Get(); err != nil {
		log.Println("获取用户信息失败,用户id:", userid, ",失败原因:", err.Error())
		return common.BackServerError(ctx, 206)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 207, "没有找到该用户")
	}

	// 重置该用户密码
	oldSalt := userModel.Salt
	oldPassword := userModel.Password

	userModel.Salt = tools.RandomString()
	userModel.Password = tools.Md5([]byte(data.Password + userModel.Salt))
	if _, err := userModel.Update("salt,password"); err != nil {
		log.Printf("重置用户密码失败,用户信息:%#v,错误原因:%s\n", userModel, err.Error())
		return common.BackServerError(ctx, 207)
	}

	// 删除该用户的所有重置token
	if err := find.Delete(userModel.UserID); err != nil {
		log.Printf("删除用户的重置token失败,用户id:%d,错误原因:%s\n", userModel, err.Error())

		userModel.Salt = oldSalt
		userModel.Password = oldPassword
		if _, err := userModel.Update("salt,password"); err != nil {
			log.Printf("还原用户原有salt和密码失败,用户信息%#v,错误原因:%s\n", userModel, err.Error())
		}

		return common.BackServerError(ctx, 208)
	}

	return common.BackOk(ctx)
}
