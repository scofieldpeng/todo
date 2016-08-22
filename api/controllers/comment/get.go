package comment

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/todo/api/libs/common"
	"github.com/scofieldpeng/todo/api/models/comment"
	"github.com/scofieldpeng/todo/api/models/todo"
	"log"
	"net/http"
	"strconv"
)

// Comments 获取评论列表
func List(ctx echo.Context) error {
	todoid, err := strconv.Atoi(ctx.Param("todoid"))
	if err != nil || todoid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 201, "请传入正确的todoid")
	}
	useridInterface := ctx.Get("userid")
	userid, ok := useridInterface.(int)
	if !ok || userid < 1 {
		log.Printf("从系统中获取userid失败,失败原因:%s\n", err.Error())
		return common.BackServerError(ctx, 202)
	}

	// 获取todo详情
	todoModel := todo.New()
	todoModel.ID = todoid
	if exsit, err := todoModel.Get(); err != nil {
		log.Printf("从数据库获取todo信息失败,todoid:%d,详情:%s\n", todoid, err.Error())
		return common.BackServerError(ctx, 203)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 204, "该todo不存在")
	}

	if todoModel.UserID != userid {
		return common.BackError(ctx, http.StatusBadRequest, 205, "没有权限查看该todo")
	}

	// 获取该todo的所有评论信息
	commentModel := comment.New()
	commentModel.TodoID = todoid
	commentList, err := commentModel.List()
	if err != nil {
		log.Println("从数据库获取评论列表数据失败,查询todoid:", todoid, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 206)
	}

	return ctx.JSON(http.StatusOK, commentList)
}
