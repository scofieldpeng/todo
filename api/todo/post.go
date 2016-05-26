package todo

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/models/todo"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "time"
    "log"
    "strconv"
)

// Insert 插入一条数据
func Insert(ctx echo.Context) error {
    userid,ok := ctx.Get("userid").(int)
    if !ok {
        userid = 0
    }
    if userid == 0 {
        return common.BackServerError(ctx,201)
    }

    getUserid,err := strconv.Atoi(ctx.Param("userid"))
    if err != nil {
        getUserid = 0
    }
    if userid != getUserid && userid != 0 {
        return common.BackError(ctx,http.StatusBadRequest,202,"请传入正确的userid")
    }

    var todoData todo.Todo
    if err := common.GetBodyStruct(ctx,&todoData);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,203,"请上传正确的数据")
    }
    if todoData.UserID != userid && userid != 0{
        if userid == 0 {
            log.Println("从已登录状态中转化userid结果为0")
        }
        return common.BackError(ctx,http.StatusBadRequest,204,"您没有权限新建,可能需要登录")
    }
    todoData.ID = 0
    todoData.UserID = userid
    if len(todoData.TodoName) < 1 || len(todoData.TodoName) > 300 {
        return common.BackError(ctx,http.StatusBadRequest,205,"todo的名称长度不合法,长度区间请在1-300个字符之间")
    }
    todoData.CategoryID = 0
    todoData.CreateTime = int(time.Now().Unix())
    todoData.Status = StatusDefault

    if _,err := todoData.Insert();err != nil {
        log.Println("插入TODO数据失败,失败原因:",err)
        return common.BackServerError(ctx,206)
    }

    return ctx.JSON(http.StatusOK,todoData)

}
