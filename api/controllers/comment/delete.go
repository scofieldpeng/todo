package comment

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/todo/api/libs/common"
	"github.com/scofieldpeng/todo/api/models/comment"
	"log"
	"net/http"
	"strconv"
)

// Delete 删除某条评论
func Delete(ctx echo.Context) error {
	todoid, err := strconv.Atoi(ctx.Param("todoid"))
	if err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 201, "请输入正确的todoid")
	}
	commentid, err := strconv.Atoi(ctx.Param("commentid"))
	if err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 202, "倾诉如正确的commentid")
	}

	useridInterface := ctx.Get("userid")
	userid, ok := useridInterface.(int)
	if !ok {
		return common.BackServerError(ctx, 203)
	}

	// 检查该条评论是否存在
	commentModel := comment.New()
	commentModel.TodoID = todoid
	commentModel.CommentID = commentid
	if exsit, err := commentModel.Get(); err != nil {
		log.Printf("获取评论详情失败,todoid:%d,评论id:%d,错误原因:%s\n", todoid, commentid, err.Error())
		return common.BackServerError(ctx, 204)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 205, "评论不存在")
	}

	if commentModel.UserID != userid {
		return common.BackError(ctx, http.StatusBadRequest, 206, "没有权限操作")
	}

	deleteModel := comment.New()
	deleteModel.CommentID = commentid

	if _, err := deleteModel.Delete(); err != nil {
		log.Println("删除评论失败,要删除的评论id:", commentid, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 207)
	}

	return common.BackOk(ctx)
}
