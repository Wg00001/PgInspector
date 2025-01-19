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
	DB      map[name]*DBConfig
	Table   map[name]*TableConfig
	Task    map[name]*TaskConfig
	Log     map[level]*LogConfig
	Alert   map[level]*AlertConfig
}

type name string
type level int

type DefaultConfig struct {
	DefaultDriver     string
	DefaultLogLevel   level
	DefaultAlertLevel level
}

type DBConfig struct {
	DBName name
	Driver string
	DSN    string
}

type TableConfig struct {
	TableName name
	DBConfig  string
}

type TaskConfig struct {
	TaskName   name
	CheckCycle time.Duration
	Todo       []*inspector.Inspect
}

type LogConfig struct {
	LogLevel level
	Table    *TableConfig
}

type AlertConfig struct {
	AlertLevel level
}
