package usecase

import (
	"PgInspector/entities/config"
	"log"
	"sync"
)

/**
 * @description: configs usecase
 * @author Wg
 * @date 2025/2/4
 */

var (
	globalConfig config.Config
	mu           sync.RWMutex
)

func GetConfig() config.Config {
	mu.RLock()
	defer mu.RUnlock()
	return globalConfig
}

func InitConfig() {
	mu.Lock()
	defer mu.Unlock()
	globalConfig = config.Config{
		Default: config.DefaultConfig{},
		DB:      make(map[config.Name]config.DBConfig),
		Table:   make(map[config.Name]config.TableConfig),
		Task:    make(map[config.Name]config.TaskConfig),
		Log:     make(map[config.Level]config.LogConfig),
		Alert:   make(map[config.Level]config.AlertConfig),
	}
}

func AddConfigs[T config.DefaultConfig | config.DBConfig | config.TableConfig | config.TaskConfig | config.LogConfig | config.AlertConfig](configs ...T) {
	if configs == nil || len(configs) == 0 {
		log.Println("AddConfigs params is nil or empty")
		return
	}
	rangeFunc := func(f func(cfg T)) {
		for _, v := range configs {
			f(v)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	switch t := any(configs[0]).(type) {
	case config.DefaultConfig:
		rangeFunc(func(cfg T) {
			globalConfig.Default = any(cfg).(config.DefaultConfig)
		})
	case config.TableConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.TableConfig)
			globalConfig.Table[val.TableName] = val
		})
	case config.DBConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.DBConfig)
			globalConfig.DB[val.DBName] = val
		})
	case config.TaskConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.TaskConfig)
			globalConfig.Task[val.TaskName] = val
		})
	case config.LogConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.LogConfig)
			globalConfig.Log[val.LogLevel] = val
		})
	case config.AlertConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.AlertConfig)
			globalConfig.Alert[val.AlertLevel] = val
		})
	default:
		log.Printf("type of config nonsupport to Add: %s\n", t)
	}
}
