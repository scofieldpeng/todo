package todo

import (
    "github.com/labstack/echo"
    "strconv"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
    "github.com/scofieldpeng/todo/models/user"
)

// Update 更新一条数据
func Update(ctx echo.Context) error {
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

    var updateTODO todo.Todo
    if err := common.GetBodyStruct(ctx,&updateTODO);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,204,"请上传正确的数据")
    }
    if updateTODO.Star < 1 {
        updateTODO.Star = 1
    }
    if updateTODO.Star > 4 {
        updateTODO.Star = 4
    }
    if updateTODO.Score < 1 {
        updateTODO.Score = 1
    }

    // 检查原有数据是否存在
    todoModel := todo.New()
    todoModel.ID = todoid
    if exist,err := todoModel.Get();err != nil {
        log.Println("获取todo数据失败,tododid:",todoid,",失败原因:",err)
        return common.BackServerError(ctx,205)
    } else if !exist {
        return common.BackError(ctx,http.StatusBadRequest,206,"没有找到该TODO")
    }
    if todoModel.UserID != userid {
        return common.BackError(ctx,http.StatusBadRequest,207,"没有找到该TODO")
    }

    updateTODO.UserID = userid
    updateTODO.ID = todoid
    if _,err := updateTODO.UpdateByID();err != nil {
        log.Printf("更新todo数据时出现错误,要更新的数据:%#v,错误原因:%s\n",updateTODO,err)
        return common.BackServerError(ctx,208)
    }

    // 如果将状态更新为已完成,将用户的总积分增加,未完成数量减1
    // 如果将状态由完成更新为在做,将用户总积分减少,未完成数量加1
    if updateTODO.Status == StatusFinish && todoModel.Status != StatusFinish{
        userModel := user.New()
        userModel.UserID = userid
        if _,err := userModel.Incr(todoModel.Score,"score");err != nil {
            log.Printf("更新用户的积分值失败,用户id:%d,要递增的数量:%d,错误原因:%#v\n",todoModel.UserID,todoModel.Score,err.Error())
        }
        if _,err := userModel.Decr(1,"unfinish_num");err != nil {
            log.Printf("减少用户的未完成数量失败,用户id:%d,错误原因:%#v\n",todoModel.UserID,err.Error())
        }
    }
    if todoModel.Status == StatusFinish && updateTODO.Status != StatusFinish {
        userModel := user.New()
        userModel.UserID = userid
        if _,err := userModel.Decr(todoModel.Score,"score");err != nil {
            log.Printf("更新用户的积分值失败,用户id:%d,错误原因:%#v\n",todoModel.UserID,err.Error())
        }
        if _,err := userModel.Incr(1,"unfinish_num");err != nil {
            log.Printf("增加用户的积分值失败,用户id:%d,错误原因:%#v\n",todoModel.UserID,err.Error())
        }
    }

    return ctx.JSON(http.StatusOK,updateTODO)
}

// TODO RegularUpdate 更新定期TODO
func RegularUpdate(ctx echo.Context) error {
    return nil
}