package user

import (
	"errors"
	"fmt"
	"github.com/scofieldpeng/mysql-go"
)

// User 用户结构体
type User struct {
	UserID      int    `json:"userid"       xorm:"not null INT(10) pk autoincr 'userid'"` // 用户id
	Email       string `json:"email"        xorm:"not null VARCHAR(100) index"`           // 邮箱
	UserName    string `json:"user_name"    xorm:"not null INT(10) index"`                // 用户登录名
	Password    string `json:"-"            xorm:"not null VARCHAR(32)"`                  // 密码
	Salt        string `json:"-"            xorm:"not null VARCHAR(8)"`                   // SALT
	CreateTime  int    `json:"create_time"  xorm:"not null INT(10)"`                      // 账号创建时间
	LastLogin   int    `json:"last_login"   xorm:"not null INT(11)"`                      // 最近一次登录时间
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
	if u.UserID == 0 && u.UserName == "" && u.Email == ""{
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

// IncrScore 增加某用户某条目的数量,传入要递增的积分值和该条目在数据库的字段名
func (u *User) Incr(incrNum int, incrFieldName string) (int64, error) {
	return mysql.Select().XormEngine().Id(u.UserID).Incr(incrFieldName, incrNum).Update(u)
}

// DecrScore 减少某用户某条目的数量,传入要减去的积分值和该条目在数据库的字段名
func (u *User) Decr(decrNum int, decrFieldName string) (int64, error) {
	return mysql.Select().XormEngine().Id(u.UserID).Decr(decrFieldName, decrNum).Where(fmt.Sprintf("%s>?", decrFieldName), decrNum).Update(u)
}
