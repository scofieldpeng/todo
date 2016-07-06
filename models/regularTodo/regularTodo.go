package regulartodo

import(
	"github.com/scofieldpeng/mysql-go"
)

// RegularTodo 定期todo模型
type RegularTodo struct {
	RegularTodoID int    `json:"regular_todoid" xorm:"not null INT(11) pk autoincr 'regular_todoid'"` // 主键
	TodoName      string `json:"todo_name"      xorm:"not null INT(11)"`                              // todo名称
	Userid        int    `json:"userid"         xorm:"not null INT(11) index"`                        // 用户id
	Type          int    `json:"type"           xorm:"not null INTYINT(1) default '1'"`               // 周期类型,1每天,2每周,3每月
	RegularDate   string `json:"regular_date"   xorm:"not null VARCHAR(100)"`                         // 周期时间,2为每周哪几天,3为每月哪几天
	CreateTime    int    `json:"create_time"    xorm:"not null INT(11) default '0'"`                  // 创建时间
	UpdateTime    int    `json:"update_time"    xorm:"not null INT(11) default '0'"`                  // 更新时间
	CategoryID    int    `json:"categoryid"     xorm:"not null INT(11) default '0' 'categoryid'"`     // todo分类
	StartTime     int    `json:"start_time"     xorm:"not null INT(11) default '0'"`                  // 开始时间
	EndTime       int    `json:"end_time"       xorm:"not null INT(11) default '0'"`                  // 结束时间
	Star          int    `json:"star"           xorm:"not null TINYINT(1) default '1'"`               // 重要程度
	Score         int    `json:"score"          xorm:"not null INT(11) default '1'"`                  // 分数
	Status        int    `json:"status"         xorm:"not null TINYINT(1) default 0"`                 // 状态
	Remark        string `json:"remark"         xorm:"not null TEXT"`                                 // 备注
}

// New 新建一个regulartodo对象
func New() RegularTodo {
	return RegularTodo{}
}

// Insert 添加一条记录
func (rt *RegularTodo) Insert() (int64,error) {
	return mysql.Select().XormEngine().Insert(rt)
}

// Update 根据regular_todoid更新一条数据
func (rt *RegularTodo) Update(cols ...string)(int64,error) {
	engine := mysql.Select().XormEngine().NewSession()
	if len(cols) > 0 {
		engine.Cols(cols...)
	} else {
		engine.AllCols()
	}

	return engine.Id(rt.RegularTodoID).Update(rt)
}

// Get 获取一条数据
func (rt *RegularTodo) Get() (bool,error) {
	return mysql.Select().XormEngine().Get(rt)
}

// List 获取列表数据
func (rt *RegularTodo) List() ([]RegularTodo,error) {
	var list []RegularTodo
	if err := mysql.Select().XormEngine().Find(&list,rt);err != nil {
		return []RegularTodo{},err
	}

	return list,nil
}

// Page 按页面获取列表数据
func (rt *RegularTodo) Page(page,pageSize int) (int64,[]RegularTodo,error) {
	var list []RegularTodo
	totalEngine := mysql.Select().XormEngine().NewSession()
	ListEngine := mysql.Select().XormEngine().NewSession()

	total,err := totalEngine.Count(rt)
	if err != nil {
		return 0,[]RegularTodo{},err
	}
	if total < 1 {
		return 0,[]RegularTodo{},nil
	}

	if pageSize < 20 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	if err := ListEngine.Limit(pageSize,(page-1)*pageSize).Find(&list,rt);err != nil {
		return 0,[]RegularTodo{},err
	}

	return total,list,nil

}

// Delete 删除数据
func (rt *RegularTodo) Delete() (int64,error) {
	return mysql.Select().XormEngine().Delete(rt)
}