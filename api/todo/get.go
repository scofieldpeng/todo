package todo

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
    "strconv"
    "github.com/scofieldpeng/todo/models/comment"
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

// Detail 获取详情
func Detail(ctx echo.Context) error {
    userid,err := strconv.Atoi(ctx.Param("userid"))
    if err != nil || userid < 1{
        return common.BackError(ctx,http.StatusBadRequest,201,"用户id不正确")
    }
    todoid,_ := strconv.Atoi(ctx.Param("todoid"))
    if err != nil {
        return common.BackError(ctx,http.StatusBadRequest,202,"请输入合法的todoid")
    }

    todoModel := todo.New()
    todoModel.UserID = userid
    todoModel.ID = todoid
    if exsit,err := todoModel.Get();err != nil {
        log.Println("获取todo详情失败,获取的userid为:",userid,",错误原因:",err)
        return common.BackServerError(ctx,203)
    } else if !exsit {
        return common.BackError(ctx,http.StatusBadRequest,204,"没有找到该todo")
    }

    return ctx.JSON(http.StatusOK,todoModel)
}

// Comments 获取评论列表
func Comments(ctx echo.Context) error {
    todoid,err := strconv.Atoi(ctx.Param("todoid"))
    if err != nil || todoid < 1 {
        return common.BackError(ctx,http.StatusBadRequest,201,"请传入正确的todoid")
    }
    useridInterface := ctx.Get("userid")
    userid,ok := useridInterface.(int)
    if !ok || userid < 1{
        log.Printf("从系统中获取userid失败,失败原因:%s\n",err.Error())
        return common.BackServerError(ctx,202)
    }

    // 获取todo详情
    todoModel := todo.New()
    todoModel.ID = todoid
    if exsit,err := todoModel.Get();err != nil {
        log.Printf("从数据库获取todo信息失败,todoid:%d,详情:%s\n",todoid,err.Error())
        return common.BackServerError(ctx,203)
    } else if !exsit {
        return common.BackError(ctx,http.StatusBadRequest,204,"该todo不存在")
    }

    if todoModel.UserID != userid {
        return common.BackError(ctx,http.StatusBadRequest,205,"没有权限查看该todo")
    }

    // 获取该todo的所有评论信息
    commentModel := comment.New()
    commentModel.TodoID = todoid
    commentList,err := commentModel.List()
    if err != nil {
        log.Println("从数据库获取评论列表数据失败,查询todoid:",todoid,",错误原因:",err.Error())
        return common.BackServerError(ctx,206)
    }

    return ctx.JSON(http.StatusOK,commentList)
}