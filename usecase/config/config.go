package config

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"fmt"
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
		Ai:      config.AgentConfig{},
		AiTask:  make(map[config.Name]*config.AgentTaskConfig),
		KBase:   make(map[config.Name]*config.KnowledgeBaseConfig),
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

type ParamType interface {
	config.DefaultConfig | config.DBConfig | config.TaskConfig | config.LogConfig | config.AlertConfig |
		config.AgentConfig | config.AgentTaskConfig | config.KnowledgeBaseConfig | *insp.Tree
}

type GetType interface {
	*config.DefaultConfig | *config.DBConfig | *config.TaskConfig | *config.LogConfig | *config.AlertConfig |
		*config.AgentConfig | *config.AgentTaskConfig | *config.KnowledgeBaseConfig | *insp.Tree
}

func Sets[T ParamType](configs ...T) (err error) {
	for i := range configs {
		err = Set(configs[i])
		if err != nil {
			return err
		}
	}
	return
}

func Set[T ParamType](cfg T) error {
	mu.Lock()
	defer mu.Unlock()
	switch t := any(cfg).(type) {
	case config.DefaultConfig:
		Config.Default = t
	case config.DBConfig:
		Config.DB[t.Name] = &t
	case config.TaskConfig:
		Config.Task[t.Name] = &t
	case config.LogConfig:
		Config.Log[t.ID] = &t
	case config.AlertConfig:
		Config.Alert[t.ID] = &t
	case config.AgentConfig:
		Config.Ai = t
	case *insp.Tree:
		Insp = t
	case config.AgentTaskConfig:
		Config.AiTask[t.Name] = &t
	case config.KnowledgeBaseConfig:
		Config.KBase[t.Name] = &t
	default:
		return fmt.Errorf("type of config nonsupport to Add: %s\n", t)
	}
	return nil
}

func Del[T ParamType](cfg T) error {
	mu.Lock()
	defer mu.Unlock()
	switch t := any(cfg).(type) {
	case config.DefaultConfig:
		Config.Default = config.DefaultConfig{} // 非 map 类型保持清空值
	case config.DBConfig:
		delete(Config.DB, t.Name)
	case config.TaskConfig:
		delete(Config.Task, t.Name)
	case config.LogConfig:
		delete(Config.Log, t.ID)
	case config.AlertConfig:
		delete(Config.Alert, t.ID)
	case config.AgentConfig:
		Config.Ai = config.AgentConfig{}
	case *insp.Tree:
		Insp = nil // 指针类型置空
	case config.AgentTaskConfig:
		delete(Config.AiTask, t.Name)
	case config.KnowledgeBaseConfig:
		delete(Config.KBase, t.Name)
	default:
		return fmt.Errorf("type of config nonsupport to Del: %T", t) // 修正错误提示类型格式符
	}
	return nil
}

func Get[T GetType](target T) (res T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("config get fail: params = %#v", res)
		}
	}()
	mu.RLock()
	defer mu.RUnlock()

	switch t := any(target).(type) {
	case *config.DefaultConfig:
		res = any(Config.Default).(T)
	case *config.DBConfig:
		if db, ok := Config.DB[t.Name]; ok {
			res = any(db).(T)
		} else {
			err = fmt.Errorf("DB config %q not found", t.Name)
		}
	case *config.TaskConfig:
		if task, ok := Config.Task[t.Name]; ok {
			res = any(task).(T)
		} else {
			err = fmt.Errorf("task config %q not found", t.Name)
		}
	case *config.LogConfig:
		if log, ok := Config.Log[t.ID]; ok {
			res = any(log).(T)
		} else {
			err = fmt.Errorf("log config %q not found", t.ID)
		}
	case *config.AlertConfig:
		if alert, ok := Config.Alert[t.ID]; ok {
			res = any(alert).(T)
		} else {
			err = fmt.Errorf("alert config %q not found", t.ID)
		}
	case *config.AgentConfig:
		res = any(Config.Ai).(T)
	case *insp.Tree:
		res = any(Insp).(T) // 直接返回指针
	case *config.AgentTaskConfig:
		if task, ok := Config.AiTask[t.Name]; ok {
			res = any(task).(T)
		} else {
			err = fmt.Errorf("agent task %q not found", t.Name)
		}
	case *config.KnowledgeBaseConfig:
		if kb, ok := Config.KBase[t.Name]; ok {
			res = any(kb).(T)
		} else {
			err = fmt.Errorf("knowledge base %q not found", t.Name)
		}
	default:
		err = fmt.Errorf("unsupported config type: %T", t)
	}

	return
}
