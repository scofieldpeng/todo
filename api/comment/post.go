package comment

import (
    "github.com/labstack/echo"
    "github.com/scofieldpeng/todo/models/comment"
    "github.com/scofieldpeng/todo/libs/common"
    "net/http"
    "strconv"
    "github.com/scofieldpeng/todo/models/todo"
    "log"
    "strings"
    "time"
)

// Post 新建一条comment
func Post(ctx echo.Context) error {
    var newComment comment.Comment
    if err := common.GetBodyStruct(ctx,&newComment);err != nil {
        return common.BackError(ctx,http.StatusBadRequest,201,"请输入正确的数据")
    }

    todoid,err := strconv.Atoi(ctx.Param("todoid"))
    if err != nil || todoid < 1 {
        return common.BackError(ctx,http.StatusBadRequest,202,"请输入正确的todoid")
    }
    useridInterface := ctx.Get("userid")
    userid,ok := useridInterface.(int)
    if !ok {
        return common.BackError(ctx,http.StatusBadRequest,203,"没有权限")
    }

    // 检查tooid信息是否合法
    newComment.CommentID = 0
    if todoid != newComment.TodoID {
        return common.BackError(ctx,http.StatusBadRequest,204,"todoid不正确")
    }
    if newComment.UserID != userid {
        return common.BackError(ctx,http.StatusBadRequest,205,"userid不正确")
    }
    newComment.Content = strings.TrimSpace(newComment.Content)
    if newComment.Content == "" {
        return common.BackError(ctx,http.StatusBadRequest,206,"评论内容不能为空")
    }

    // 检查todo是否存在
    todoModel := todo.New()
    todoModel.ID = todoid
    if exsit,err := todoModel.Get();err != nil{
        log.Println("获取todo详情失败,todoid:",todoid,",错误原因:",err.Error())
        return common.BackServerError(ctx,207)
    } else if !exsit {
        return common.BackError(ctx,http.StatusBadRequest,208,"该todo不存在")
    }

    // 检查该todo是否和评论人id一致
    if todoModel.ID != todoid {
        return common.BackError(ctx,http.StatusBadRequest,209,"无权操作该todo")
    }

    // 如果有父级commentid,检查是否存在该parentcommentid
    if newComment.ParentCommentID > 0 {
        parentCommentModel := comment.New()
        parentCommentModel.CommentID = newComment.ParentCommentID
        if exsit,err := parentCommentModel.Get();err != nil {
            log.Println("获取parentCommentid详情失败,parentcommentid:",newComment.ParentCommentID,",错误原因:",err.Error())
            return common.BackServerError(ctx,210)
        } else if !exsit {
            return common.BackError(ctx,http.StatusBadRequest,211,"parent_commentid不存在")
        }
    }

    newComment.Time = int(time.Now().Unix())
    if _,err := newComment.Insert();err != nil {
        log.Printf("插入comment数据失败!插入comment数据:%#v,错误原因:%s\n",newComment,err.Error())
        return common.BackServerError(ctx,212)
    }

    return ctx.JSON(http.StatusOK,newComment)
}
