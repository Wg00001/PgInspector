package usecase

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"log"
	"sync"
)

/**
 * @description: 配置中心
 * @author Wg
 * @date 2025/2/4
 */

var (
	Config = config.Config{
		Default: config.DefaultConfig{},
		Task:    make(map[config.Name]*config.TaskConfig),
		DB:      make(map[config.Name]*config.DBConfig),
		Log:     make(map[config.ID]*config.LogConfig),
		Alert:   make(map[config.ID]*config.AlertConfig),
		Ai:      config.AiConfig{},
		AiTask:  make(map[config.Name]*config.AiTaskConfig),
	}
	Insp = insp.NewTree()
	mu   sync.RWMutex
)

func RLock() {
	mu.RLock()
}

func RUnlock() {
	mu.RUnlock()
}

func GetConfig() *config.Config {
	mu.RLock()
	defer mu.RUnlock()
	return &Config
}

func GetDbConfig(name config.Name) *config.DBConfig {
	mu.RLock()
	defer mu.RUnlock()
	if res, ok := Config.DB[name]; ok {
		return res
	}
	return nil
}

func GetInsp(path config.Name) *insp.Node {
	mu.RLock()
	defer mu.RUnlock()
	return Insp.GetNode(path.Str())
}

func GetAllInsp() []*insp.Node {
	mu.RLock()
	defer mu.RUnlock()
	return Insp.AllInsp
}

func GetTaskConfig(name config.Name) *config.TaskConfig {
	mu.RLock()
	defer mu.RUnlock()
	if res, ok := Config.Task[name]; ok {
		return res
	}
	return nil
}

func GetLoggerConfig(id config.ID) *config.LogConfig {
	mu.RLock()
	defer mu.RUnlock()
	if res, ok := Config.Log[id]; ok {
		return res
	}
	return nil
}

func AddConfigs[T config.DefaultConfig | config.DBConfig | config.TaskConfig | config.LogConfig | config.AlertConfig | config.AiConfig | *insp.Tree | config.AiTaskConfig](configs ...T) {
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
			Config.DB[val.Name] = &val
		})
	case config.TaskConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.TaskConfig)
			Config.Task[val.TaskName] = &val
		})
	case config.LogConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.LogConfig)
			Config.Log[val.LogID] = &val
		})
	case config.AlertConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.AlertConfig)
			Config.Alert[val.AlertID] = &val
		})
	case config.AiConfig:
		rangeFunc(func(cfg T) {
			Config.Ai = any(cfg).(config.AiConfig)
		})
	case *insp.Tree:
		rangeFunc(func(cfg T) {
			val := any(cfg).(*insp.Tree)
			Insp = val
		})
	case config.AiTaskConfig:
		rangeFunc(func(cfg T) {
			val := any(cfg).(config.AiTaskConfig)
			Config.AiTask[val.AiTaskName] = &val
		})
	default:
		log.Printf("type of config_adapter nonsupport to Add: %s\n", t)
	}
}
