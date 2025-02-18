package logger

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
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
	Log(InspLog, insp.Result)
	GetID() config.ID
}

type InspLog struct {
	Timestamp time.Time
	TaskName  config.Name
	DBName    config.Name
	TaskID    string //task批次编号
	Result    string
}

func (l InspLog) WithErr(err error) InspLog {
	l.Result = err.Error()
	return l
}

func (l InspLog) WithJSON(val any) InspLog {
	marshal, err := json.Marshal(val)
	if err != nil {
		return l.WithErr(err)
	}
	l.Result = string(marshal)
	return l
}
