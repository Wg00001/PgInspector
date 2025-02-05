package config

import (
	"PgInspector/entities/inspector"
	"time"
)

/**
 * @description: 配置的实体定义
 * @author Wg
 * @date 2025/1/19
 */

type Config struct {
	Default DefaultConfig
	DB      map[Name]DBConfig
	Table   map[Name]TableConfig
	Task    map[Name]TaskConfig
	Log     map[Level]LogConfig
	Alert   map[Level]AlertConfig
}

type Name string
type Level int

type DefaultConfig struct {
	DefaultDriver     string
	DefaultLogLevel   Level
	DefaultAlertLevel Level
}

type DBConfig struct {
	DBName Name
	Driver string
	DSN    string
}

type TableConfig struct {
	TableName Name
	DBConfig  string
}

type TaskConfig struct {
	TaskName   Name
	Async      bool
	CheckCycle time.Duration
	Todo       []*inspector.Inspect
}

type LogConfig struct {
	LogLevel Level
	Table    *TableConfig
}

type AlertConfig struct {
	AlertLevel Level
}
