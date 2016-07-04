package comment

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/models/comment"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "strconv"
    "log"
    "strings"
    "time"
)

// Put 更新评论
// 评论只能更新内容,其他无法更新
func Put(ctx echo.Context) error {
    var updateComment comment.Comment
    if err := common.GetBodyStruct(ctx,&updateComment);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,201,"请输入正确的数据")
    }

    todoid,err := strconv.Atoi(ctx.Param("todoid"))
    if err != nil || todoid < 1 {
        return common.BackError(ctx,http.StatusBadRequest,202,"请输入正确的todoid")
    }

    commentid,err := strconv.Atoi(ctx.Param("commentid"))
    if err != nil || commentid < 1 {
        return common.BackError(ctx,http.StatusBadRequest,203,"请输入正确的commentid")
    }
    useridInterface := ctx.Get("userid")
    userid, ok := useridInterface.(int)
    if !ok {
        return common.BackServerError(ctx,204)
    }

    // 检查原有评论id是否存在
    previousComment := comment.New()
    previousComment.CommentID = commentid
    if exsit,err := previousComment.Get();err != nil {
        log.Println("获取评论id数据失败,要获取的评论id:",commentid,",错误原因:",err.Error())
        return common.BackServerError(ctx,205)
    } else if !exsit {
        return common.BackError(ctx,http.StatusBadRequest,206,"该评论不存在")
    }
    if previousComment.UserID != userid {
        return common.BackError(ctx,http.StatusBadRequest,207,"无权操作该评论")
    }

    // update
    previousComment.Content = strings.TrimSpace(updateComment.Content)
    if previousComment.Content == "" {
        return common.BackError(ctx,http.StatusBadRequest,208,"评论内容不能为空")
    }
    previousComment.Time = int(time.Now().Unix())

    if _,err := previousComment.UpdateByCommentID("content");err != nil {
        log.Printf("更新评论失败,更新内容%#v,失败原因:%s\n",previousComment,err.Error())
        return common.BackServerError(ctx,209)
    }

    return ctx.JSON(http.StatusBadRequest,previousComment)
}
