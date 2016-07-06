package todo

import (
	"github.com/labstack/echo"
	"github.com/scofieldpeng/todo/libs/common"
	"github.com/scofieldpeng/todo/models/todo"
	"log"
	"net/http"
	"strconv"
	"github.com/scofieldpeng/todo/models/regulartodo"
)

// List 获取列表数据
func List(ctx echo.Context) error {
	userid, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || userid < 1 {
		return common.BackError(ctx, http.StatusBadRequest, 201, "用户id不正确")
	}

	todoModel := todo.New()
	todoModel.UserID = userid
	list, err := todoModel.List()
	if err != nil {
		log.Println("获取todo列表数据失败,获取的userid为:", userid, ",错误原因:", err)
		return common.BackServerError(ctx, 202)
	}
	total := len(list)
	if total == 0 {
		list = make([]todo.Todo, 0)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
		"list":  list,
	})
}

// RegularList 获取定期todo列表
func RegularList(ctx echo.Context) error {
	useridInterface := ctx.Get("userid")
	userid,ok := useridInterface.(int)
	if !ok {
		return common.BackServerError(ctx,201)
	}

	page,err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize,err := strconv.Atoi(ctx.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	categoryid,err := strconv.Atoi(ctx.QueryParam("categoryid"))
	if err != nil || categoryid < 1 {
		categoryid = 0
	}
	startTime,err := strconv.Atoi(ctx.QueryParam("start_time"))
	if err != nil || startTime < 1 {
		startTime = 0
	}
	endTime,err := strconv.Atoi(ctx.QueryParam("end_time"))
	if err != nil || endTime < 1 {
		endTime = 0
	}
	if endTime < startTime && (startTime != 0 && endTime != 0) {
		return common.BackError(ctx,http.StatusBadRequest,202,"结束时间不能小于开始时间")
	}
	createTime,err := strconv.Atoi(ctx.QueryParam("create_time"))
	if err != nil || createTime < 1{
		createTime  = 0
	}
	status,err := strconv.Atoi(ctx.QueryParam("status"))
	if err != nil || ( status < StatusDefault || status > StatusPaused ) {
		status = StatusDefault
	}
	todoType,err := strconv.Atoi(ctx.QueryParam("type"))
	if err != nil || (todoType < OnceTodo || todoType > MonthlyTodo) {
		todoType = OnceTodo
	}
	star,err := strconv.Atoi(ctx.QueryParam("star"))
	if err != nil || (star < NormalTodo || star > EmergencyImportantTodo) {
		star = NormalTodo
	}


	regularTodoModel := regulartodo.New()
	regularTodoModel.CategoryID = categoryid
	regularTodoModel.Userid = userid
	regularTodoModel.StartTime = startTime
	regularTodoModel.EndTime = endTime
	regularTodoModel.CreateTime = createTime
	regularTodoModel.Type = todoType
	regularTodoModel.Star = star
	regularTodoModel.Status = status

	total,list,err := regularTodoModel.Page(page,pageSize)
	if err != nil {
		log.Printf("获取定期todo列表数据失败,regularTodo数据:%#v,错误原因:%s\n",regularTodoModel,err.Error())
		return common.BackServerError(ctx,203)
	}

	return ctx.JSON(http.StatusOK,map[string]interface{}{
		"total":total,
		"list":list,
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
	userid,ok := useridInterface.(int)
	if !ok {
		return common.BackServerError(ctx,201)
	}
	inputUserid,err := strconv.Atoi(ctx.Param("userid"))
	if err != nil || inputUserid < 1 {
		return common.BackError(ctx,http.StatusBadRequest,202,"请输入正确的用户名")
	}
	if userid != inputUserid {
		return common.BackError(ctx,http.StatusBadRequest,203,"授权不通过")
	}

	todoid,err := strconv.Atoi(ctx.Param("todoid"))
	if err != nil || todoid < 1 {
		return common.BackError(ctx,http.StatusBadRequest,204,"todoid不正确")
	}

	regularTodoModel := regulartodo.New()
	regularTodoModel.RegularTodoID = todoid
	if exsit,err := regularTodoModel.Get();err != nil {
		log.Println("获取regular_todo详情失败,regular_todoid:",todoid,",错误原因:",err.Error())
		return common.BackServerError(ctx,205)
	} else if !exsit {
		return common.BackError(ctx,http.StatusBadRequest,205,"该todo不存在")
	}

	return ctx.JSON(http.StatusOK,regularTodoModel)
}
