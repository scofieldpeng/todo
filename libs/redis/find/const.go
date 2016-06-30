package find

import "fmt"

const (
	RedisToken      string = "findtoken:string:token:%s"         // redis中token对应的userid的key名
	RedisUserTokens string = "findtoken:sets:user:%d"            // redis中某个用户所拥有的token名称列表
	TokenLastTime          = "findtoken:string:user:%d:lasttime" // redis中某个用户最近一次请求找回密码时间戳

	RedisTokenExpire = 604800 // TOKEN相关用户过期时间
	TokenTimeExpire  = 60     // TOKEN最近一次请求时间过期时间
)

// generateRedisToken 生成redistoken值
func generateRedisToken(token string) string {
	return fmt.Sprintf(RedisToken,token)
}