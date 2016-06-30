package email

// Email 类型
type Email struct {
	to      []string          // 收件人
	tpl     string            // 模板
	tplVals map[string]string // 模板变量
}

// SetTo 设置单个发件人
func (e *Email) SetTo(to string) *Email {
	e.to = append(e.to, to)
	return e
}

// SetToMany 设置多个发件人
func (e *Email) SetToMany(to []string) *Email {
	e.to = append(e.to, to...)
	return e
}

// Settpl 设置tpl
func (e *Email) SetTpl(tpl string) *Email {
	e.tpl = tpl
	return e
}

// SetTplVal 设置模板变量
func (e *Email) SetTplVal(tpl,val string) *Email {
	e.tplVals[tpl] = val
	return e
}

// SetTplVals 设置多个模板的变量
func (e *Email) SetTplVals(tpls map[string]string) *Email {
	for key,val := range tpls {
		e.tplVals[key] = val
	}
	return e
}

// Send 发送邮件,如果出错,返回error
func (e *Email) Send() error {
	return nil
}


func New() Email {
	return Email{
		to:      make([]string, 0),
		tplVals: make(map[string]string),
	}
}
