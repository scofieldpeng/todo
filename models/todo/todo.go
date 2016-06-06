package todo

import (
	"errors"
	"github.com/scofieldpeng/mysql-go"
)

// TODO模型
type Todo struct {
	ID         int    `json:"id"          xorm:"not null INT(11) pk autoincr 'id'"`              // 主键id
	TodoName   string `json:"todo_name"   xorm:"not null VARCHAR(50)"`                           // 标题
	UserID     int    `json:"userid"      xorm:"not null INT(10) index(user_todo) 'userid'"`     // 用户id
	CategoryID int    `json:"category_id" xorm:"not null INT(10) index(user_todo) 'categoryid'"` // 分类id
	CreateTime int    `json:"create_time" xorm:"not null INT(10)"`                               // 创建时间
	StartTime  int    `json:"start_time"  xorm:"not null INT(10)"`                               // 开始时间
	EndTime    int    `json:"end_time"    xorm:"not null INT(10)"`                               // 结束时间
	Status     int    `json:"status"      xorm:"not null TINYINT(1) default 0 index(user_todo)"` // 状态,0待做,1再做,2完成,3放弃
	Remark     string `json:"remark"        xorm:"not null TEXT"`                                // 备注
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

// List 获取列表数据
func (t *Todo) List() ([]Todo, error) {
	var list []Todo
	if err := mysql.Select().XormEngine().Find(&list, t); err != nil {
		return []Todo{}, err
	}

	return list, nil
}
