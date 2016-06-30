package common

import (
	"github.com/labstack/echo"
	"net/http"
)

// BackError 向客户端返回错误,httpCode为http错误码,errcode为返回的app错误码,errmsg为返回的APP错误信息
func BackError(ctx echo.Context, httpCode, errcode int, errmsg string) error {
	return ctx.JSON(httpCode, map[string]interface{}{
		"errcode": errcode,
		"errmsg":  errmsg,
	})
}

// BackUnauthorized 需要登录授权
func BackUnAuthorized(ctx echo.Context) error {
	return BackError(ctx, http.StatusUnauthorized, 401, "unauthorized error")
}

// BackServerError 服务器内部错误
func BackServerError(ctx echo.Context, errcode int) error {
	return BackError(ctx, http.StatusInternalServerError, errcode, "server error")
}

// BackOk 通用返回
func BackOk(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]bool{
		"status": true,
	})
}
