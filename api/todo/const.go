package todo

const (
	StatusDefault = 0 + iota // todo状态(未开始)
	StatusDoing              // todo状态(进行中)
	StatusFinish             // todo状态(已完成)
	StatusPaused             // todo状态(暂停中)
)

const (
	OnceTodo    = 0 + iota // 单次todo
	DailyTodo              // 每日todo
	WeeklyTodo             // 每周todo
	MonthlyTodo            // 每月todo
)

const (
	NormalTodo             = 0 + iota // 一般todo
	ImportantTodo                     // 重要todo
	EmergencyTodo                     // 紧急todo
	EmergencyImportantTodo            // 紧急且重要todo
)
