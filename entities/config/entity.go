package config

import (
	"PgInspector/entities/insp"
	"time"
)

/**
 * @description: 配置的实体定义
 * @author Wg
 * @date 2025/1/19
 */

type Config struct {
	Default DefaultConfig
	Task    map[Name]*TaskConfig
	DB      map[Name]*DBConfig
	Log     map[ID]*LogConfig
	Alert   map[ID]*AlertConfig
	Insp    *insp.Tree
}

type Name string
type ID int

type DefaultConfig struct {
	DefaultDriver     string
	DefaultLogLevel   ID
	DefaultAlertLevel ID
}

type DBConfig struct {
	Name   Name
	Driver string
	DSN    string
}

type LogConfig struct {
	LogID  ID
	Driver string
	Header map[string]string
}

type AlertConfig struct {
	AlertLevel ID
	//todo:alert
}

type TaskConfig struct {
	TaskName     Name
	Time         *Cron
	AllInspector bool
	//Priority     int
	//Async        bool
	//todo:定时任务

	TargetDB   []Name
	LogLevel   ID
	AlertLevel ID

	Todo    []Name
	NotTodo []Name
}

type Cron struct {
	Duration time.Duration
	AtTime   []string
	Weekly   []time.Weekday
	Monthly  []int
}
