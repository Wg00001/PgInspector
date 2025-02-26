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

	Ai     AiConfig
	AiTask map[Name]*AiTaskConfig
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
	LogID  ID
	Driver string
	Header map[string]string
}

type AlertConfig struct {
	AlertID ID
	Driver  string
	Header  map[string]string
}

// ---task(任务)相关配置

type TaskConfig struct {
	TaskName     Name
	Cron         *Cron
	AllInspector bool
	//todo:async

	LogID    ID
	TargetDB []Name

	Todo    []Name
	NotTodo []Name
}

type Cron struct {
	Duration time.Duration
	AtTime   []string
	Weekly   []time.Weekday
	Monthly  []int
}

// ---Ai相关配置

// AiConfig 用户只能指定一个全局Ai，所有的分析均由此Ai完成
type AiConfig struct {
	//AiName      Name
	Driver      string
	Api         string
	ApiKey      string
	Model       string
	Temperature float64
}

type AiTaskConfig struct {
	AiTaskName Name
	Cron       *Cron
	LogID      ID
	LogFilter  LogFilter
	AlertID    ID
}

type LogFilter struct {
	// 时间范围：Timestamp 需在 [StartTime, EndTime] 之间
	StartTime time.Time
	EndTime   time.Time
	TaskNames []Name   // TaskName 匹配列表
	DBNames   []Name   // DBName 匹配列表
	TaskIDs   []string // TaskID 匹配列表
	InspNames []Name   // Insp匹配列表
}
