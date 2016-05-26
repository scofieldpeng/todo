package user

import (
	"errors"
	"github.com/scofieldpeng/mysql-go"
)

// User 用户结构体
type User struct {
	UserID      int    `json:"userid" xorm:"not null INT(10) pk autoincr 'userid'"` // 用户id
	UserName    string `json:"user_name" xorm:"not null INT(10) index"`             // 用户登录名
	Password    string `json:"-" xorm:"not null VARCHAR(32)"`                       // 密码
	Salt        string `json:"-" xorm:"not null VARCHAR(8)"`                        // SALT
	CreateTime  int    `json:"create_time" xorm:"not null INT(10)"`                 // 账号创建时间
	LastLogin   int    `json:"last_login" xorm:"not null INT(11)"`                  // 最近一次登录时间
	UnfinishNum int    `json:"unfinish_num" xorm:"not null INT(11)"`                // 未完成数量
}

// New 新建一个用户结构体对象
func New() User {
	return User{}
}

// Insert 插入一条数据
func (u *User) Insert() (int64, error) {
	return mysql.Select().XormEngine().Insert(u)
}

// Get 获取用户名
func (u *User) Get() (bool, error) {
	if u.UserID == 0 && u.UserName == "" {
		return false, errors.New("获取用户信息必须指定用户id或者用户名")
	}
	return mysql.Select().XormEngine().Get(u)
}

// Update 更新数据
func (u *User) Update(cols ...string) (int64, error) {
	engine := mysql.Select().XormEngine().NewSession()
	if len(cols) == 0 {
		engine.AllCols()
	} else {
		engine.Cols(cols...)
	}

	if u.UserID != 0 {
		engine.Id(u.UserID)
	}
	if u.UserName != "" {
		engine.Where("user_name=?", u.UserName)
	}

	return engine.Update(u)
}
