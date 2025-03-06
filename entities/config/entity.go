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
	Log     map[ID]*LogConfig
	Alert   map[ID]*AlertConfig

	Ai     AgentConfig
	AiTask map[Name]*AgentTaskConfig
	KBase  map[Name]*KnowledgeBaseConfig
	//Insp    *insp.Tree //insp不放在此处，避免循环引用
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
	ID     ID
	Driver string
	Header map[string]string
}

type AlertConfig struct {
	ID     ID
	Driver string
	Header map[string]string
}

// ---task(任务)相关配置

type TaskConfig struct {
	Name         Name
	Cron         *Cron
	AllInspector bool
	//todo:async

	LogID    ID
	TargetDB []Name

	Todo    []Name
	NotTodo []Name
}

type Cron struct {
	CronTab  string
	Duration time.Duration
	AtTime   []string
	Weekly   []time.Weekday
	Monthly  []int
}

// ---Ai Agent 相关配置

// AgentConfig 用户只能指定一个全局Ai，所有的分析均由此Ai完成
type AgentConfig struct {
	//AiName      Name
	Driver        string
	Url           string
	ApiKey        string
	Model         string
	Temperature   float64
	SystemMessage string
}

type AgentTaskConfig struct {
	Name          Name
	Cron          *Cron
	LogID         ID
	LogFilter     LogFilter
	AlertID       ID
	KBase         []Name
	KBaseTopN     int
	KBaseMaxLen   int
	SystemMessage string
}

type LogFilter struct {
	// 时间范围：Timestamp 需在 [StartTime, EndTime] 之间
	StartTime time.Time
	EndTime   time.Time
	TaskNames []Name   // Name 匹配列表
	DBNames   []Name   // DBName 匹配列表
	TaskIDs   []string // TaskID 匹配列表
	InspNames []Name   // Insp匹配列表
}

type KnowledgeBaseConfig struct {
	Name
	Driver string
	Value  map[string]string
}
