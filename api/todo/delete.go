package todo

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "strconv"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
    "github.com/scofieldpeng/todo/models/user"
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

// TODO DeleteRegular 删除定期todo
func RegularDelete(ctx echo.Context) error {
    return nil
}