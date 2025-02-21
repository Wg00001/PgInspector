package logger

import (
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"encoding/json"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

//根据config来设置多个logger，不同实现的logger的配置格式不同，在usecase中转换成标准格式送到adapter中各自解析。
//用户配置task时可以指定loggerId，没有指定则使用0号。

type Logger interface {
	Log(Content, db.Result)
	GetID() config.ID
	Init(cfg *config.LogConfig) (Logger, error)
}

type Content struct {
	Timestamp time.Time
	TaskName  config.Name
	DBName    config.Name
	TaskID    string //task批次编号
	Result    string
}

func (l Content) WithErr(err error) Content {
	l.Result = err.Error()
	return l
}

func (l Content) WithJSON(val any) Content {
	marshal, err := json.Marshal(val)
	if err != nil {
		return l.WithErr(err)
	}
	l.Result = string(marshal)
	return l
}
