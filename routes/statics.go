package routes

import (
	"github.com/labstack/echo/middleware"
	"github.com/scofieldpeng/todo/libs/common"
	"github.com/scofieldpeng/todo/controllers/index"
)

func init() {
	common.Echo.Pre(middleware.RemoveTrailingSlash())

    // statics访问
	common.Echo.Static("/statics","statics")
	// index页面
	common.Echo.Get("/",index.Index)
}
