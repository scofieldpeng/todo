package todo

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "strconv"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
    "github.com/scofieldpeng/todo/models/user"
    "github.com/scofieldpeng/todo/models/regulartodo"
)

// Delete 删除数据
func Delete(ctx echo.Context) error {
    userid,ok := ctx.Get("userid").(int)
    if !ok {
        userid = 0
    }
    if userid == 0 {
        return common.BackServerError(ctx,201)
    }

    getUserid,err := strconv.Atoi(ctx.Param("userid"))
    if err != nil || getUserid < 1 || getUserid != userid{
        return common.BackError(ctx,http.StatusBadRequest,202,"请输入正确的userid")
    }
    todoid,err := strconv.Atoi(ctx.Param("todoid"))
    if err != nil || todoid< 1 {
        return common.BackError(ctx,http.StatusBadRequest,203,"请输入正确的todoid")
    }

    // 检查原有数据是否存在
    todoModel := todo.New()
    todoModel.ID = todoid
    if exist,err := todoModel.Get();err != nil {
        log.Println("获取todo数据失败,tododid:",todoid,",失败原因:",err)
        return common.BackServerError(ctx,204)
    } else if !exist {
        return common.BackError(ctx,http.StatusBadRequest,205,"没有找到该TODO")
    }
    if todoModel.UserID != userid {
        return common.BackError(ctx,http.StatusBadRequest,206,"没有找到该TODO")
    }

    if _,err := todoModel.Delete();err != nil {
        log.Println("删除todo数据时失败,失败原因:",err.Error())
        return common.BackServerError(ctx,207)
    }

    // 如果用户当前todo未完成,删除后要减去该用户的已完成数量,单位1
    if todoModel.Status != StatusFinish {
        userModel := user.New()
        userModel.UserID = userid
        if _,err := userModel.Decr(1,"unfinish_num");err != nil {
            log.Printf("减少用户的未完成数量失败,用户id:%d,错误原因:%#v\n",todoModel.UserID,err.Error())
        }
    }

    return ctx.JSON(http.StatusOK,map[string]bool{
        "status":true,
    })
}

// deleteRegular 删除定期todo
func RegularDelete(ctx echo.Context) error {
    useridInterface := ctx.Get("userid")
    userid,ok := useridInterface.(int)
    if !ok {
        return common.BackServerError(ctx,201)
    }
    inputUserid,err := strconv.Atoi(ctx.Param("userid"))
    if err != nil || inputUserid < 1 {
        return common.BackError(ctx,http.StatusBadRequest,202,"请输入正确的用户名")
    }
    if userid != inputUserid {
        return common.BackError(ctx,http.StatusBadRequest,203,"授权不通过")
    }

    todoid,err := strconv.Atoi(ctx.Param("todoid"))
    if err != nil || todoid < 1 {
        return common.BackError(ctx,http.StatusBadRequest,204,"todoid不正确")
    }

    regularTodoModel := regulartodo.New()
    regularTodoModel.RegularTodoID = todoid
    if exsit,err := regularTodoModel.Get();err != nil {
        log.Println("获取regular_todo详情失败,regular_todoid:",todoid,",错误原因:",err.Error())
        return common.BackServerError(ctx,205)
    } else if !exsit {
        return common.BackError(ctx,http.StatusBadRequest,205,"该todo不存在")
    }
    if regularTodoModel.Userid != userid {
        return common.BackError(ctx,http.StatusBadRequest,206,"没有权限删除该todo")
    }

    deleteTodoModel := regularTodoModel
    deleteTodoModel.RegularTodoID = todoid
    if _,err := deleteTodoModel.Delete();err != nil {
        log.Println("删除regular_todo失败,要删除的reguarl_todoid:",todoid,",错误原因:",err.Error())
        return common.BackServerError(ctx,207)
    }

    // 将所属于该regulartodo的todo(未完成)一并删除
    todoModel := todo.New()
    todoModel.RegularTodoID = todoid
    todoModel.Status = StatusDefault
    deleteNum,err := todoModel.Delete()
    if err != nil {
        log.Printf("删除regular_todo生成的待做todo列表数据失败,要删除的todoid:%d,错误原因:%s\n",todoid,err.Error())
        if _,err := regularTodoModel.Insert();err != nil {
            log.Printf("删除regular_todo生成的待做todo列表数据失败,回滚regulartodo失败,回滚regular_todo数据:%#v,错误原因:%s\n",regularTodoModel,err.Error())
            return common.BackServerError(ctx,208)
        }
        return common.BackServerError(ctx,209)
    }
    // 将用户未完成数量减少
    userModel := user.New()
    userModel.UserID = userid
    if _,err := userModel.Decr(int(deleteNum),"unfinish_num");err != nil {
        log.Println("减小用户未完成数量失败,用户id:",userid,",错误原因:",err.Error())
        return common.BackServerError(ctx,210)
    }

    return common.BackOk(ctx)
}