package config

import "time"

type Config struct {
	GlobalConfig
	PgConfigs []PgConfig
	LogDB     DbConfig
}

type PgConfig struct {
	DbConfig
	InspectorConfig
	AlertConfig
}

type DbConfig struct {
	Name     string `yaml:"name"`     // 数据库标识
	Host     string `yaml:"host"`     // 数据库地址
	Port     int    `yaml:"port"`     // 数据库端口
	Username string `yaml:"username"` // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	Dbname   string `yaml:"dbname"`   // 数据库名
	Charset  string `yaml:"charset"`  // 字符集
}

type InspectorConfig struct {
	CheckInterval time.Duration `yaml:"check_interval"` // 巡检间隔
	AlertEmail    []string      `yaml:"alert_email"`    // 报警通知的邮箱列表
}

type AlertConfig struct {
	QueryTimeout   time.Duration `yaml:"query_timeout"`   // 查询超时阈值
	MaxConnections int           `yaml:"max_connections"` // 最大连接数
	LogLevel       string        `yaml:"log_level"`       // 日志级别
	LogFile        string        `yaml:"log_file"`        // 日志文件路径
}

type GlobalConfig struct {
	LogLevel        int  `yaml:"log_level"`        // 全局日志级别
	AlertEnabled    bool `yaml:"alert_enabled"`    // 是否启用报警
	SaveToDB        bool `yaml:"save_to_db"`       // 是否保存巡检结果到数据库
	ResultRetention int  `yaml:"result_retention"` // 巡检结果保留天数
}
