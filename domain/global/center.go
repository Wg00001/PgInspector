package global

import (
	"PgInspector/domain/task"
	"database/sql"
)

var (
	dbPool    []*sql.DB
	tablePool []string
	taskPool  []task.Task //todo:循环引用了
)
