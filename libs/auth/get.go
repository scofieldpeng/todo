package auth

import (
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	libRedis "github.com/scofieldpeng/redis-go"
	"log"
)

// GetTokenFromCookie 从cookie中获取token值
func GetTokenFromCookie(ctx echo.Context) string {
	cookie, err := ctx.Request().Cookie(Token_Cookie_Name)
	if err != nil {
		return ""
	}

	return cookie.Value()
}

// GetUseridFromRedis 根据token从redis中获取用户id
func GetUseridFromRedis(token string) int {
	if token == "" {
		return 0
	}

	redisConn := libRedis.Pool("default").Get()
	defer redisConn.Close()

	userid, err := redis.Int(redisConn.Do("GET", RedisLoginTokenName(token)))
	if err != nil {
		if err != redis.ErrNil {
			log.Println("从redis中获取信息失败,错误原因:", err.Error())
		}
	}

	return userid
}
