package logger

import "time"

// Log 每个task对应一个log,使用PostgreSQL的TimescaleDB保存
type Log struct {
	//主键,时序数据库
	Time time.Time
	//任务名,对应的任务
	Task string
	//具体的监控指标和对应数值
	//使用PostgreSQL的JSONB类型来存储
	Metric map[string]float32
}
