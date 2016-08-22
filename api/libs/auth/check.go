package auth

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/todo/api/libs/common"
	"fmt"
)

// Check 用户api接口检查用户授权的middleware
func Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// 获取本次请求接口是否需要授权
		if !ApiNeedAuth(ctx.Path(), RequestMethod(ctx.Request().Method())) {
			return next(ctx)
		}

		// 需要授权,检查该用户是否成功登陆
		cookie := GetTokenFromCookie(ctx)
		if cookie == "" {
			fmt.Println("need auth,but not found cookie")
			return common.BackUnAuthorized(ctx)
		}
		fmt.Println("auth cookie:",cookie)

		userid := GetUseridFromRedis(cookie)
		if userid == 0 {
			fmt.Println("need auth,but the cookie value is 0")
			return common.BackUnAuthorized(ctx)
		}

		ctx.Set("userid", userid)
		return next(ctx)

	}
}
