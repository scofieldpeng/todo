package todo

import (
    "github.com/labstack/echo"
    "strconv"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
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

    return ctx.JSON(http.StatusOK,updateTODO)
}