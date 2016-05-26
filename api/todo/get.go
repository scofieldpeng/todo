package todo

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
    "strconv"
)



// List 获取列表数据
func List(ctx echo.Context) error {
    userid,err := strconv.Atoi(ctx.Param("userid"))
    if err != nil || userid < 1{
        return common.BackError(ctx,http.StatusBadRequest,201,"用户id不正确")
    }

    todoModel := todo.New()
    todoModel.UserID = userid
    list,err := todoModel.List()
    if err != nil {
        log.Println("获取todo列表数据失败,获取的userid为:",userid,",错误原因:",err)
        return common.BackServerError(ctx,202)
    }
    total := len(list)
    if  total == 0 {
        list = make([]todo.Todo,0)
    }

    return ctx.JSON(http.StatusOK,map[string]interface{}{
        "total":total,
        "list":list,
    })
}