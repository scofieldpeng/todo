package auth

import "fmt"

const (
	Token_Cookie_Name   string = "todo_logintoken" // cookie中用户登录token名称
	Token_Cookie_Expire int    = 252000            // cookie中用户登录token的过期时间
)

// RedisLoginTokenName 获取redis中用户token的字段名称
func RedisLoginTokenName(token string) string {
	return fmt.Sprintf("todo:logintoken:%s", token)
}
