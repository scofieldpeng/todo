package routes

import (
	"github.com/scofieldpeng/todo/api/todo"
	"github.com/scofieldpeng/todo/api/user"
	"github.com/scofieldpeng/todo/libs/auth"
	"github.com/scofieldpeng/todo/libs/common"
)

func init() {
	api := common.Echo.Group("/api", auth.Check)
	apiV1 := api.Group("/v1")

	// TODO列表
	apiV1.GET(`/todo/:userid`, todo.List)
	// 插入TODO
	apiV1.POST(`/todo/:userid`, todo.Insert)
    // TODO详情
	apiV1.GET(`/todo/:userid/:todoid`, todo.Detail)
	// TODO更新
	apiV1.PUT(`/todo/:userid/:todoid`, todo.Update)
	// TODO删除
	apiV1.DELETE(`/todo/:userid/:todoid`, todo.Delete)

	// TODO 评论
	apiV1.Get(`/todo/:todoid/comments`,todo.Comments) // todo评论列表
	

	apiV1.POST(`/login`, user.Login)       // 登录用户
	apiV1.POST(`/register`, user.Register) // 注册用户
	apiV1.POST(`/find`, user.Find)         // 找回密码申请接口
	apiV1.POST(`/reset`, user.ResetPwd)    // 重置密码
}
