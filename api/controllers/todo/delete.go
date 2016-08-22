package todo

import (
	"github.com/labstack/echo"
	apiComment "github.com/scofieldpeng/todo/api/controllers/comment"
	"github.com/scofieldpeng/todo/api/libs/common"
	"github.com/scofieldpeng/todo/api/models/comment"
	"github.com/scofieldpeng/todo/api/models/regulartodo"
	"github.com/scofieldpeng/todo/api/models/todo"
	"github.com/scofieldpeng/todo/api/models/user"
	"log"
	"net/http"
	"strconv"
)

// Delete 删除数据
func Delete(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		userid = 0
	}
	if userid == 0 {
		return common.BackServerError(ctx, 201)
	}

	getUserid, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || getUserid < 1 || getUserid != userid {
		return common.BackError(ctx, http.StatusBadRequest, 202, "请输入正确的userid")
	}
	todoid, err := strconv.Atoi(ctx.Param("todoid"))
	if err != nil || todoid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 203, "请输入正确的todoid")
	}

	// 检查原有数据是否存在
	todoModel := todo.New()
	todoModel.ID = todoid
	if exist, err := todoModel.Get(); err != nil {
		log.Println("获取todo数据失败,tododid:", todoid, ",失败原因:", err)
		return common.BackServerError(ctx, 204)
	} else if !exist {
		return common.BackError(ctx, http.StatusBadRequest, 205, "没有找到该TODO")
	}
	if todoModel.UserID != userid {
		return common.BackError(ctx, http.StatusBadRequest, 206, "没有找到该TODO")
	}

	if _, err := todoModel.Delete(); err != nil {
		log.Println("删除todo数据时失败,失败原因:", err.Error())
		return common.BackServerError(ctx, 207)
	}

	// 如果用户当前todo未完成,删除后要减去该用户的已完成数量,单位1
	if todoModel.Status != StatusFinish {
		userModel := user.New()
		userModel.UserID = userid
		if _, err := userModel.Decr(1, "unfinish_num"); err != nil {
			log.Printf("减少用户的未完成数量失败,用户id:%d,错误原因:%#v\n", todoModel.UserID, err.Error())
		}
	}

	return ctx.JSON(http.StatusOK, map[string]bool{
		"status": true,
	})
}

// deleteRegular 删除定期todo
func RegularDelete(ctx echo.Context) error {
	useridInterface := ctx.Get("userid")
	userid, ok := useridInterface.(int)
	if !ok {
		return common.BackServerError(ctx, 201)
	}
	inputUserid, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || inputUserid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 202, "请输入正确的用户名")
	}
	if userid != inputUserid {
		return common.BackError(ctx, http.StatusBadRequest, 203, "授权不通过")
	}

	regularTodoid, err := strconv.Atoi(ctx.Param("todoid"))
	if err != nil || regularTodoid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 204, "todoid不正确")
	}

	regularTodoModel := regulartodo.New()
	regularTodoModel.RegularTodoID = regularTodoid
	if exsit, err := regularTodoModel.Get(); err != nil {
		log.Println("获取regular_todo详情失败,regular_todoid:", regularTodoid, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 205)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 206, "该todo不存在")
	}
	if regularTodoModel.Userid != userid {
		return common.BackError(ctx, http.StatusBadRequest, 207, "没有权限删除该todo")
	}

	deleteTodoModel := regularTodoModel
	deleteTodoModel.RegularTodoID = regularTodoid
	if _, err := deleteTodoModel.Delete(); err != nil {
		log.Println("删除regular_todo失败,要删除的reguarl_todoid:", regularTodoid, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 208)
	}

	// 将所属于该regulartodo的todo一并删除
	todoModel := todo.New()
	todoModel.RegularTodoID = regularTodoid
	todoList, err := todoModel.List()
	if err != nil {
		log.Println("获取regular_todo的todo失败,regular_todoid:", regularTodoid, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 209)

	}
	todoids := make([]int, 0)
	for _, todo := range todoList {
		todoids = append(todoids, todo.ID)
	}

	//
	if len(todoList) > 0 {
		if _, err := todoModel.DeleteByIDs(todoids); err != nil {
			log.Printf("删除regular_todo生成的待做todo列表数据失败,要删除的todoid:%d,错误原因:%s\n", regularTodoid, err.Error())
			if _, err := regularTodoModel.Insert(); err != nil {
				log.Printf("删除regular_todo生成的待做todo列表数据失败,回滚regulartodo失败,回滚regular_todo数据:%#v,错误原因:%s\n", regularTodoModel, err.Error())
				return common.BackServerError(ctx, 210)
			}
			return common.BackServerError(ctx, 211)
		}

		// 将regular_todo下属的相关的comment删除
		commentModel := comment.New()
		commentModel.Type = apiComment.Normal_TODO
		if _, err := commentModel.DeleteByTodoIDs(todoids); err != nil {
			log.Printf("删除regular_todo下属的一次性todo评论数据失败,删除的todoids:%#v,错误原因:%s\n", todoids, err.Error())
			// TODO rollback
			return common.BackServerError(ctx, 213)
		}
	}

	// 删除该regular_todo的评论数据
	commentModel := comment.New()
	commentModel.Type = apiComment.Rugular_TODO
	commentModel.TodoID = regularTodoid
	if _, err := commentModel.Delete(); err != nil {
		log.Println("删除regular_todo的评论数据失败,删除的regular_todoid:", regularTodoid, ",错误原因:", err.Error())
		// TODO rollback
		return common.BackServerError(ctx, 214)
	}
	return common.BackOk(ctx)
}
