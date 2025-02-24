package logger

import (
	"PgInspector/entities/config"
	"time"
)

/**
 * @description: logger ReadLog函数的filter，使用functional option模式
 * @author Wg
 * @date 2025/2/23
 */

type Filter struct {
	// 时间范围：Timestamp 需在 [StartTime, EndTime] 之间
	StartTime time.Time
	EndTime   time.Time

	TaskNames []config.Name // TaskName 匹配列表
	DBNames   []config.Name // DBName 匹配列表
	TaskIDs   []string      // TaskID 匹配列表
}

// NewFilter 新建一个过滤器，默认获取前一周的所有数据
func NewFilter(opts ...Option) Filter {
	filter := &Filter{StartTime: time.Now().AddDate(0, 0, -7)}
	for _, opt := range opts {
		opt(filter)
	}
	return *filter
}

type Option func(*Filter)

func WithTimeRange(start, end time.Time) Option {
	return func(filter *Filter) {
		filter.StartTime = start
		filter.EndTime = end
	}
}

func WithStartTime(start time.Time) Option {
	return func(filter *Filter) {
		filter.StartTime = start
	}
}

func WithEndTime(end time.Time) Option {
	return func(filter *Filter) {
		filter.EndTime = end
	}
}

func WithTaskNames(names ...config.Name) Option {
	return func(o *Filter) {
		o.TaskNames = names
	}
}

func WithDBNames(names ...config.Name) Option {
	return func(o *Filter) {
		o.DBNames = names
	}
}

func WithTaskIDs(ids ...string) Option {
	return func(o *Filter) {
		o.TaskIDs = ids
	}
}
