package common

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	Debug bool // 是否是debug模式
	Echo  *echo.Echo
)

func init() {
	Echo = echo.New()
	Echo.Use(middleware.Gzip())
	Echo.Use(middleware.Logger())
	Echo.Use()
}
