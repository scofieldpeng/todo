package auth

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/redis-go"
    "time"
	"net"
	"github.com/scofieldpeng/go-libs/tools"
	"strconv"
)


// createToken 生成token,首先尝试网卡mac+当前时间戳(微秒)+用户id然后md5,如果没有找到可用的网卡,直接当前时间戳(微秒) + 用户id然后md5
func createToken(userid int) string {
interfaces,err := net.Interfaces()
if len(interfaces) > 0 && err == nil{
	return tools.Md5([]byte(string(interfaces[0].HardwareAddr) + strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(userid)))
}

return tools.Md5([]byte( strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(userid) ))
}

// SaveToken 保存用户token,将会生成一个随机token,写入到redis和cookie中,返回的第一个参数为生成的token
func SaveToken(ctx echo.Context,userid int) (string,error) {
	token := createToken(userid)
	if err := SaveToRedis(userid,token);err != nil {
		return "",err
	}

	SaveToCookie(ctx,token)
	return token,nil
}

// SaveToken 将token写入到redis中,传入用户id和要写入的token即可,如果出错,将会返回error
func SaveToRedis(userid int, token string) error {
	redisConn := redis.Pool().Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("SETEX", RedisLoginTokenName(token), Token_Cookie_Expire+1, userid); err != nil {
		return err
	}

	return nil
}

// SaveToCookie 将用户的token值存储到cookie中
func SaveToCookie(ctx echo.Context, token string) {
    cookie := new(echo.Cookie)
    cookie.SetName(Token_Cookie_Name)
    cookie.SetValue(token)
    cookie.SetExpires(time.Unix(time.Now().Unix() + int64(Token_Cookie_Expire),0))
    cookie.SetPath("/")

    ctx.SetCookie(cookie)
}
