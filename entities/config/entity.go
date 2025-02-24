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
	//Insp    *insp.Tree
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

type TaskConfig struct {
	TaskName     Name
	Time         *Cron
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

type AiConfig struct {
	AiName      Name
	Api         string
	ApiKey      string
	Model       string
	Temperature float64
}
