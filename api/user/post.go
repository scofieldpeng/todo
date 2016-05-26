package user

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "regexp"
    "errors"
    "github.com/scofieldpeng/todo/models/user"
    "log"
    "github.com/scofieldpeng/go-libs/tools"
    "github.com/scofieldpeng/todo/libs/auth"
    "time"
)

// loginData 登录请求数据
type loginData struct {
    UserName string `json:"user_name"`
    Password string `json:"password"`
}

// isValidUserName 检查是否是有效的用户名
func validUserName(username string) error {
    if matched,err := regexp.Match("^[a-zA-z_-]{6,20}$",[]byte(username));err != nil || !matched{
        return errors.New("非法用户名")
    }

    return nil
}

// checkPassword 检查密码是否正确
func validPassword(password string) error {
    if matched,err := regexp.Match("^[a-zA-z0-9]{1,32}$",[]byte(password));err != nil || !matched{
        return errors.New("非法用户名")
    }
    return nil
}

// Login 用户登录
func Login(ctx echo.Context) error {
    var logindata loginData
    if err := common.GetBodyStruct(ctx,&logindata);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,201,"请输入正确的数据")
    }

    if err := validUserName(logindata.UserName);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,202,err.Error())
    }

    if err := validPassword(logindata.Password);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,203,err.Error())
    }

    // 检查用户名是否存在
    userModel := user.New()
    userModel.UserName = logindata.UserName
    if exsit,err := userModel.Get();err != nil {
        log.Println("获取用户数据失败!用户名:",logindata.UserName,",错误原因:",err.Error())
        return common.BackServerError(ctx,204)
    } else if !exsit {
        return common.BackError(ctx,http.StatusBadRequest,205,"用户不存在")
    }

    // 检查用户密码
    if userModel.Password != tools.Md5([]byte(logindata.Password + userModel.Salt)) {
        return common.BackError(ctx,http.StatusBadRequest,206,"用户密码不正确")
    }

    if _,err := auth.SaveToken(ctx,userModel.UserID);err != nil {
        log.Println("生成login token时失败!用户名:",userModel.UserName,",错误原因:",err.Error())
        return common.BackServerError(ctx,207)
    }

    lastLogin := userModel.LastLogin

    // 更新用户上次登录时间
    userModel.LastLogin = int(time.Now().Unix())
    if _,err := userModel.Update();err != nil {
        log.Println("更新用户最近一次登录时间时出错,错误原因:",err.Error())
    }
    userModel.LastLogin = lastLogin

    return ctx.JSON(http.StatusOK,userModel)
}
