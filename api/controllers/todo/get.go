package todo

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/todo/api/libs/common"
	"github.com/scofieldpeng/todo/api/models/regulartodo"
	"github.com/scofieldpeng/todo/api/models/todo"
	"log"
	"net/http"
	"strconv"
)

// List 获取分页列表数据
func List(ctx echo.Context) error {
	userid, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || userid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 201, "用户id不正确")
	}
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(ctx.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	categoryid,err := strconv.Atoi(ctx.QueryParam("categoryid"))
	if err != nil || categoryid < 1{
		categoryid = 0
	}
	order := ctx.QueryParam("order")
	if err != nil || (order != StarOrder && order != StartTimeOrder || order != EndTimeOrder) {
		order = DefaultOrder
	}

	var lc todo.ListCondtion
	startTime,err := strconv.Atoi(ctx.QueryParam("start_time"))
	if err != nil || startTime < 0{
		lc.StartTime = startTime
	}
	endTime,err := strconv.Atoi(ctx.QueryParam("end_time"))
	if err != nil {
		lc.EndTime = endTime
	}
	startCreateTime,err := strconv.Atoi(ctx.QueryParam("start_create_time"))
	if err != nil || startCreateTime < 0 {
		lc.StartCreateTime = 0
	}
	endCreateTime,err := strconv.Atoi(ctx.QueryParam("end_create_time"))
	if err != nil || endCreateTime < 0 {
		lc.EndCreateTime = 0
	}
	star,err := strconv.Atoi(ctx.QueryParam("star"))
	if err != nil || star < NormalTodo || star > EmergencyImportantTodo && star != -1{
		lc.Star = NormalTodo
	}
	status, err := strconv.Atoi(ctx.QueryParam("status"))
	if err != nil || (status < StatusDefault || status > StatusPaused) && status != -1{
		status = StatusDefault
	}

	todoModel := todo.New()
	todoModel.UserID = userid
	if status != -1 {
		todoModel.Status = status
	}
	todoModel.CategoryID = categoryid
	total, list, err := todoModel.Page(page, pageSize, lc,order)
	if err != nil {
		log.Println("获取todo列表数据失败,获取的userid为:", userid, ",错误原因:", err)
		return common.BackServerError(ctx, 202)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
		"list":  list,
	})
}

// RegularList 获取定期todo列表
func RegularList(ctx echo.Context) error {
	useridInterface := ctx.Get("userid")
	userid, ok := useridInterface.(int)
	if !ok {
		return common.BackServerError(ctx, 201)
	}

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(ctx.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	categoryid, err := strconv.Atoi(ctx.QueryParam("categoryid"))
	if err != nil || categoryid < 1 {
		categoryid = 0
	}
	status, err := strconv.Atoi(ctx.QueryParam("status"))
	if err != nil || (status < StatusDefault || status > StatusPaused) {
		status = StatusDefault
	}
	todoType, err := strconv.Atoi(ctx.QueryParam("type"))
	if err != nil || (todoType < OnceTodo || todoType > MonthlyTodo) {
		todoType = OnceTodo
	}
	order := ctx.QueryParam("order")
	if err != nil || (order != StarOrder && order != StartTimeOrder || order != EndTimeOrder) {
		order = DefaultOrder
	}

	var lc todo.ListCondtion
	startCreateTime,err := strconv.Atoi(ctx.QueryParam("start_create_time"))
	if err != nil || startCreateTime < 0 {
		lc.StartCreateTime = 0
	}
	endCreateTime,err := strconv.Atoi(ctx.QueryParam("end_create_time"))
	if err != nil || endCreateTime < 0 {
		lc.EndCreateTime = 0
	}
	star,err := strconv.Atoi(ctx.QueryParam("star"))
	if err != nil || star < NormalTodo || star > EmergencyImportantTodo && star != -1{
		lc.Star = 0
	}

	regularTodoModel := regulartodo.New()
	regularTodoModel.CategoryID = categoryid
	regularTodoModel.Userid = userid
	regularTodoModel.Type = todoType
	regularTodoModel.Star = star
	regularTodoModel.Status = status

	total, list, err := regularTodoModel.Page(page, pageSize,lc,order)
	if err != nil {
		log.Printf("获取定期todo列表数据失败,regularTodo数据:%#v,错误原因:%s\n", regularTodoModel, err.Error())
		return common.BackServerError(ctx, 203)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
		"list":  list,
	})
}

// Detail 获取详情
func Detail(ctx echo.Context) error {
	userid, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || userid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 201, "用户id不正确")
	}
	todoid, _ := strconv.Atoi(ctx.Param("todoid"))
	if err != nil {
		return common.BackError(ctx, http.StatusBadRequest, 202, "请输入合法的todoid")
	}

	todoModel := todo.New()
	todoModel.UserID = userid
	todoModel.ID = todoid
	if exsit, err := todoModel.Get(); err != nil {
		log.Println("获取todo详情失败,获取的userid为:", userid, ",错误原因:", err)
		return common.BackServerError(ctx, 203)
	} else if !exsit {
		return common.BackError(ctx, http.StatusBadRequest, 204, "没有找到该todo")
	}

	return ctx.JSON(http.StatusOK, todoModel)
}

// RegularDetail 获取定期todo详情
func RegularDetail(ctx echo.Context) error {
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
		return common.BackError(ctx, http.StatusBadRequest, 205, "该todo不存在")
	}
	if regularTodoModel.Userid != userid {
		return common.BackError(ctx, http.StatusBadRequest, 206, "没有权限访问该todo")
	}

	return ctx.JSON(http.StatusOK, regularTodoModel)
}
