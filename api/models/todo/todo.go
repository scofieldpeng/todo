package todo

import (
	"errors"
	"github.com/scofieldpeng/mysql-go"
)

// TODO模型
type Todo struct {
	ID            int    `json:"id"             xorm:"not null INT(11) pk autoincr 'id'"`                   // 主键id
	TodoName      string `json:"todo_name"      xorm:"not null VARCHAR(50)"`                                // 标题
	Type          int    `json:"type"           xorm:"not null TINYINT(1)"`                                 // TODO类型,0为单次,1为每日,2为每周,3为每月
	RegularTodoID int    `json:"regular_todoid" xorm:"not null INT(11) index 'regular_todoid' default '0'"` // 定期任务todoid
	UserID        int    `json:"userid"         xorm:"not null INT(10) index(user_todo) 'userid'"`          // 用户id
	CategoryID    int    `json:"category_id"    xorm:"not null INT(10) index(user_todo) 'categoryid'"`      // 分类id
	CreateTime    int    `json:"create_time"    xorm:"not null INT(10)"`                                    // 创建时间
	StartTime     int    `json:"start_time"     xorm:"not null INT(10)"`                                    // 开始时间
	EndTime       int    `json:"end_time"       xorm:"not null INT(10)"`                                    // 结束时间
	Status        int    `json:"status"         xorm:"not null TINYINT(1) default '0' index(user_todo)"`    // 状态,0待做,1再做,2完成,3放弃
	Star          int    `json:"star"           xorm:"not null TINYTIN(1) default '1'"`                     // todo的重要程度,1为一般,2重要,3紧急,4重要且紧急
	Score         int    `json:"score"          xorm:"not null INT(10) default 0"`                          // 该todo的积分,用户设置,默认为0
	Remark        string `json:"remark"         xorm:"not null TEXT"`                                       // 备注
}

// listCondition 列表查询条件
// 选择放在这里而不是api那里是因为会引起引用循坏,在没有想到更合理的办法之前先放在这里吧
type ListCondtion struct {
	StartTime       int // 开始时间
	EndTime         int // 结束时间
	StartCreateTime int // 创建时间(起)
	EndCreateTime   int // 创建时间(止)
	Star int // 重要程度
}

// New 新建一个todo结构体对象
func New() Todo {
	return Todo{}
}

// Insert 插入数据
func (t *Todo) Insert() (int64, error) {
	return mysql.Select().XormEngine().Insert(t)
}

// Update 根据ID更新一条数据
func (t *Todo) UpdateByID(cols ...string) (int64, error) {
	engine := mysql.Select().XormEngine().NewSession()
	if len(cols) == 0 {
		engine.AllCols()
	} else {
		engine.Cols(cols...)
	}

	return engine.Id(t.ID).Update(t)
}

// UpdateByRegularID 更新regularID更新
func (t *Todo) UpdateByRegularID(cols ...string)(int64,error) {
	engine := mysql.Select().XormEngine().NewSession()
	if len(cols) == 0 {
		engine.AllCols()
	} else {
		engine.Cols(cols...)
	}

	return engine.Where("regular_todoid=? AND status=?",t.RegularTodoID,t.Status).Update(t)
}

// Get 获取一条数据
func (t *Todo) Get() (bool, error) {
	return mysql.Select().XormEngine().Get(t)
}

// Delete数据
func (t *Todo) Delete() (int64, error) {
	if t.ID == 0 && t.UserID == 0 {
		return 0, errors.New("删除数据时必须指定userid或者id")
	}
	return mysql.Select().XormEngine().Delete(t)
}

// DeleteByIDs 删除一定id区间内的todo数据
func (t *Todo) DeleteByIDs(ids []int)(int64,error) {
	if len(ids) == 0 {
		return 0,nil
	}
	return mysql.Select().XormEngine().In("id",ids).Delete(t)
}

// List 获取列表数据
func (t *Todo) List(order ...string) ([]Todo, error) {
	var list []Todo
	engine := mysql.Select().XormEngine().NewSession()
	if len(order) > 0 {
		engine.OrderBy(order[0])
	}
	if err := engine.Find(&list, t); err != nil {
		return []Todo{}, err
	}

	return list, nil
}

// Page 获取分页列表数据
func (t *Todo) Page(page,pageSize int,whereCond ListCondtion,order ...string) (int64,[]Todo,error) {
	t1 := *t
	var list []Todo
	totalEngine := mysql.Select().XormEngine().NewSession()
	listEngine := mysql.Select().XormEngine()

	if whereCond.StartTime !=0 {
		totalEngine.Where("start_time>=?",whereCond.StartTime)
		listEngine.Where("start_time>=?",whereCond.StartTime)
		t1.StartTime = 0
	}
	if whereCond.EndTime != 0 {
		totalEngine.Where("end_time<=?",whereCond.EndTime)
		listEngine.Where("end_time<=?",whereCond.EndTime)
		t1.EndTime = 0
	}

	if whereCond.StartCreateTime != 0 {
		totalEngine.Where("create_time<=?",whereCond.StartCreateTime)
		listEngine.Where("create_time<=?",whereCond.EndCreateTime)
		t1.CreateTime = 0
	}
	if whereCond.Star != 0 {
		totalEngine.Where("star=?",t1.Star)
		listEngine.Where("star=?",t1.Star)
	}

	total,err := totalEngine.Count(t1)
	if err != nil {
		return 0,[]Todo{},err
	}
	if total < 1 {
		return 0,[]Todo{},nil
	}
	if len(order) > 0 {
		listEngine.OrderBy(order[0])
	}

	if pageSize < 20 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	if err := listEngine.Limit(pageSize,(page-1)*pageSize).Find(&list,t1);err != nil {
		return 0,[]Todo{},err
	}

	return total,list,nil
}
