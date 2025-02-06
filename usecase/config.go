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
	Config config.Config
	mu     sync.RWMutex
)

func GetConfig() config.Config {
	mu.RLock()
	defer mu.RUnlock()
	return Config
}

func InitConfig() {
	mu.Lock()
	defer mu.Unlock()
	Config = config.Config{
		Default: config.DefaultConfig{},
		Task:    make(map[config.Name]config.TaskConfig),
		DB:      make(map[config.Name]config.DBConfig),
		Log:     make(map[config.Level]config.LogConfig),
		Alert:   make(map[config.Level]config.AlertConfig),
		Insp:    make(map[config.Name]config.InspectConfig),
	}
}

func AddConfigs[T config.DefaultConfig | config.DBConfig | config.TaskConfig | config.LogConfig | config.AlertConfig | config.InspectConfig](configs ...T) {
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
			Config.Default = any(cfg).(config.DefaultConfig)
		})
	case config.DBConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.DBConfig)
			Config.DB[val.Name] = val
		})
	case config.TaskConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.TaskConfig)
			Config.Task[val.TaskName] = val
		})
	case config.LogConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.LogConfig)
			Config.Log[val.LogLevel] = val
		})
	case config.AlertConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.AlertConfig)
			Config.Alert[val.AlertLevel] = val
		})
	case config.InspectConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.InspectConfig)
			Config.Insp[val.InspName] = val
		})
	default:
		log.Printf("type of config_reader nonsupport to Add: %s\n", t)
	}
}
