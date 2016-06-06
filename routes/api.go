package routes

import (
    "github.com/scofieldpeng/todo/libs/common"
    "github.com/scofieldpeng/todo/libs/auth"
    "github.com/scofieldpeng/todo/api/todo"
    "github.com/scofieldpeng/todo/api/user"
)

func init() {
    api := common.Echo.Group("/api", auth.Check)
    apiV1 := api.Group("/v1")

    // TODO列表
    apiV1.GET(`/todo/:userid`, todo.List)
    // 插入TODO
    apiV1.POST(`/todo/:userid`,todo.Insert)
    apiV1.GET(`/todo/:userid/:todoid`,todo.Detail)
    // TODO更新
    apiV1.PUT(`/todo/:userid/:todoid`, todo.Update)
    // TODO删除
    apiV1.DELETE(`/todo/:userid/:todoid`, todo.Delete)

    // 用户登陆
    apiV1.POST(`/login`,user.Login)
}
