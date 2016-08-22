package todo

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/todo/api/libs/common"
	"github.com/scofieldpeng/todo/api/models/regulartodo"
	"github.com/scofieldpeng/todo/api/models/todo"
	"github.com/scofieldpeng/todo/api/models/user"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Update 更新一条数据
func Update(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok || userid < 1 {
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

	var updateTodo todo.Todo
	if err := common.GetBodyStruct(ctx, &updateTodo); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 204, "请上传正确的数据")
	}
	updateTodo.TodoName = strings.TrimSpace(updateTodo.TodoName)
	if updateTodo.TodoName == "" {
		return common.BackError(ctx, http.StatusBadRequest, 205, "todo名称不能为空")
	}
	if updateTodo.Star < NormalTodo {
		updateTodo.Star = NormalTodo
	}
	if updateTodo.Star > EmergencyImportantTodo {
		updateTodo.Star = EmergencyImportantTodo
	}
	if updateTodo.Score < 1 {
		updateTodo.Score = 1
	}

	// 检查原有数据是否存在
	todoModel := todo.New()
	todoModel.ID = todoid
	if exist, err := todoModel.Get(); err != nil {
		log.Println("获取todo数据失败,tododid:", todoid, ",失败原因:", err)
		return common.BackServerError(ctx, 205)
	} else if !exist {
		return common.BackError(ctx, http.StatusBadRequest, 206, "没有找到该TODO")
	}
	if todoModel.UserID != userid {
		return common.BackError(ctx, http.StatusBadRequest, 207, "没有找到该TODO")
	}

	updateTodo.UserID = userid
	updateTodo.ID = todoid
	if _, err := updateTodo.UpdateByID(); err != nil {
		log.Printf("更新todo数据时出现错误,要更新的数据:%#v,错误原因:%s\n", updateTodo, err)
		return common.BackServerError(ctx, 208)
	}

	// 如果将状态更新为已完成,将用户的总积分增加,未完成数量减1
	// 如果将状态由完成更新为在做,将用户总积分减少,未完成数量加1
	if updateTodo.Status == StatusFinish && todoModel.Status != StatusFinish {
		userModel := user.New()
		userModel.UserID = userid
		if _, err := userModel.Incr(todoModel.Score, "score"); err != nil {
			log.Printf("更新用户的积分值失败,用户id:%d,要递增的数量:%d,错误原因:%#v\n", todoModel.UserID, todoModel.Score, err.Error())
		}
	}
	if todoModel.Status == StatusFinish && updateTodo.Status != StatusFinish {
		userModel := user.New()
		userModel.UserID = userid
		if _, err := userModel.Decr(todoModel.Score, "score"); err != nil {
			log.Printf("更新用户的积分值失败,用户id:%d,错误原因:%#v\n", todoModel.UserID, err.Error())
		}
		if _, err := userModel.Incr(1, "unfinish_num"); err != nil {
			log.Printf("增加用户的积分值失败,用户id:%d,错误原因:%#v\n", todoModel.UserID, err.Error())
		}
	}

	return ctx.JSON(http.StatusOK, updateTodo)
}

// RegularUpdate 更新定期TODO
// 更新todo的时候需要将没有开始的属于该regular_todo的todo的信息也一并更新
func RegularUpdate(ctx echo.Context) error {
	userid, ok := ctx.Get("userid").(int)
	if !ok {
		return common.BackServerError(ctx, 201)
	}
	inputUserid, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || inputUserid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 202, "请输入正确的userid")
	}
	todoid, err := strconv.Atoi(ctx.Param("todoid"))
	if err != nil || todoid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 203, "请输入正确的todoid")
	}

	var updateTodo regulartodo.RegularTodo
	if err := common.GetBodyStruct(ctx, &updateTodo); err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 204, "请上传正确的数据")
	}
	if updateTodo.Star < NormalTodo {
		updateTodo.Star = NormalTodo
	}
	if updateTodo.Star > EmergencyImportantTodo {
		updateTodo.Star = EmergencyImportantTodo
	}
	if updateTodo.Score < 1 {
		updateTodo.Score = 1
	}
	if updateTodo.TodoName == "" {
		return common.BackError(ctx, http.StatusBadRequest, 205, "todo名称不能为空")
	}
	if updateTodo.Type < DailyTodo || updateTodo.Type > MonthlyTodo {
		return common.BackError(ctx, http.StatusBadRequest, 206, "todo_type不合法")
	}

	regularTodoModel := regulartodo.New()
	regularTodoModel.RegularTodoID = todoid
	if exsit, err := regularTodoModel.Get(); err != nil {
		log.Println("获取regular_todo失败,要获取的regular_todoid:", todoid, ",错误原因:", err.Error())
		return common.BackServerError(ctx, 206)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 207, "该todo不存在")
	}
	if regularTodoModel.Userid != userid {
		return common.BackError(ctx, http.StatusBadRequest, 208, "没有权限操作")
	}

	if _, err := updateTodo.Update(); err != nil {
		log.Printf("更新regular_todo出错,更新的todo%#v,出错原因:%s\n", updateTodo, err.Error())
		return common.BackServerError(ctx, 209)
	}

	// 更新的时候应该将所有所属该regular_todo的todo(未开始)的一并更新
	todoModel := todo.New()
	todoModel.RegularTodoID = todoid
	todoModel.TodoName = regularTodoModel.TodoName
	todoModel.Type = regularTodoModel.Type
	todoModel.CategoryID = regularTodoModel.CategoryID
	todoModel.Star = regularTodoModel.Star
	todoModel.Score = regularTodoModel.Score
	todoModel.Remark = regularTodoModel.Remark
	todoModel.Status = StatusDefault
	if _, err := todoModel.UpdateByRegularID("todo_name,type,categoryid,star,score,remark"); err != nil {
		if _, err := regularTodoModel.Update(); err != nil {
			log.Printf("更新regular_todo的todo列表失败时回滚regular_todo失败,原有的regular_todo信息:%#v,更新后的信息:%#v,错误原因:%s\n", regularTodoModel, updateTodo, err.Error())
		}
		log.Printf("更新regular_todo所属的todo失败,更新的regular_todo信息:%#v,失败原因:%s\n", updateTodo, err.Error())
	}

	return ctx.JSON(http.StatusOK, updateTodo)
}
