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

const (
	DefaultOrder   = "create_time DESC" // 默认排序
	StarOrder      = "star DESC"        // 星级排序
	StartTimeOrder = "start_time DESC"  // 创建时间排序
	EndTimeOrder   = "end_time DESC"    // 结束时间排序
)
