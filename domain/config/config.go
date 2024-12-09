package config

type Config struct {
	LogLevel        int  `yaml:"log_level"`        // 全局默认日志级别
	AlertEnabled    bool `yaml:"alert_level"`      // 全局默认报警级别
	ResultRetention int  `yaml:"result_retention"` // 巡检结果保留天数
	DBConfigs       map[dbName]*DBConfig
	TableConfigs    map[tableName]*TableConfig
	TaskConfigs     map[taskName]*TaskConfig
	LogConfigs      map[logLevel]*LogConfig
	AlertConfigs    map[alterLevel]*AlertConfig
}

type dbName string
type DBConfig struct {
	DBName     dbName `yaml:"db_name"` // 对应yaml的第二层,解析时在循环中获取
	Driver     string `yaml:"driver"`  // 数据源类型
	DSN        string `yaml:"dsn"`     // DSN
	AlterLevel int    `yaml:"alter_level"`
	LogLevel   int    `yaml:"log_level"`
	*AlertConfig
	*LogConfig
}

type tableName string
type TableConfig struct {
	TableName  tableName `yaml:"table_name"`
	DBName     string    `yaml:"db_name"`
	AlterLevel int       `yaml:"alter_level"`
	LogLevel   int       `yaml:"log_level"`
	*DBConfig
	*AlertConfig
	*LogConfig
}

type taskName string

// TaskConfig 用于控制Inspector是否执行
type TaskConfig struct {
	TaskName     taskName `yaml:"task_name"`
	CheckCycle   int      `yaml:"check_cycle"` // 巡检周期 (秒)
	AllInspector bool     `yaml:"do_all"`
	TodoList     []string `yaml:"todo"`
	NotDoList    []string `yaml:"notdo"`
}

type logLevel int

// todo
type LogConfig struct {
	LogLevel     int          `yaml:"log_level"` //主键
	LogTableName string       `yaml:"table_name"`
	LogTable     *TableConfig //将数据库作为导出
	LogFilePath  string       `yaml:"file_path"`
}

type alterLevel int

// todo
type AlertConfig struct { //
	AlterLevel int `yaml:"alter_level"`
	FeishuUrl  string
}

// todo
type GrafanaConfig struct {
	Url string
}
