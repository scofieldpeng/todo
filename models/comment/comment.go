package comment

import "github.com/scofieldpeng/mysql-go"

// Comment comment结构体对象
type Comment struct {
	CommentID       int    `json:"commentid"        xorm:"not null INT(11) pk autoincr 'commentid'"`        // 评论id
	TodoID          int    `json:"todoid"           xorm:"not null INT(11) index default '0' 'todoid'"`     // todoid
	UserID          int    `json:"userid"           xorm:"not null INT(11) default '0' 'userid"`            // userid
	ParentCommentID int    `json:"parent_commentid" xorm:"not null INT(11) default '0' 'parent_commentid'"` // 父级评论列表
	Time            int    `json:"time"             xorm:"not null INT(11) default '0'"`                    // 评论时间
	Content         string `json:"content"          xorm:"not null TEXT"`                                   // 评论内容
}

// New 新建一个comment结构体对象
func New() Comment{
    return Comment{}
}

func (c *Comment) Insert() (int64,error) {
    return mysql.Select().XormEngine().Insert(c)
}

// UpdateByCommentID 通过commentid更新数据,选择要更新的field,默认全部更新,第一个参数返回更新的数据行数,如果出错,第二个参数反馈error
func (c *Comment) UpdateByCommentID(cols ...string) (int64,error) {
    engine := mysql.Select().XormEngine()
    if len(cols) > 0 {
        engine.Cols(cols...)
    } else {
        engine.AllCols()
    }
    return engine.Id(c.CommentID).Update(c)
}

// Get 获取一条数据,第一个参数返回是否存在,如果查询失败,第二个参数返回error
func (c *Comment) Get() (bool,error) {
    return mysql.Select().XormEngine().Get(c)
}

// Delete 删除数据
func (c *Comment) Delete() (bool,error) {
    return mysql.Select().XormEngine().Delete(c)
}