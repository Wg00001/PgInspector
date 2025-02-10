package config

import (
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
	Insp    map[Name]*InspectConfig
	Log     map[Level]*LogConfig
	Alert   map[Level]*AlertConfig
}

type Name string
type Level int

type DefaultConfig struct {
	DefaultDriver     string
	DefaultLogLevel   Level
	DefaultAlertLevel Level
}

type DBConfig struct {
	Name   Name
	Driver string
	DSN    string
}

type InspectConfig struct {
	InspName Name
	SQL      string
}

type LogConfig struct {
	LogLevel  Level
	TableName Name
	//todo:log
}

type AlertConfig struct {
	AlertLevel Level
	//todo:alert
}

type TaskConfig struct {
	TaskName     Name
	Priority     int
	Async        bool
	AllInspector bool
	CheckCycle   time.Duration //todo:定时任务

	TargetDB   []Name
	LogLevel   Level
	AlertLevel Level

	Todo    []Name
	NotTodo []Name
}
