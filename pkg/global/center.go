package global

import (
	"PgInspector/pkg/task"
	"database/sql"
)

var (
	dbPool    []*sql.DB
	tablePool []string
	taskPool  []task.Task //todo:循环引用了
)
