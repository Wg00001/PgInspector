package config

import (
	"PgInspector/entities/config"
	"fmt"
	"sync"
)

/**
 * @description: 配置中心
 * @author Wg
 * @date 2025/2/4
 */

var (
	Meta  = config.ConfigMeta{}
	Index = config.ConfigIndex{
		Default:   &config.DefaultConfig{},
		Task:      make(map[config.Identity]*config.TaskConfig),
		DB:        make(map[config.Identity]*config.DBConfig),
		Log:       make(map[config.Identity]*config.LogConfig),
		Alert:     make(map[config.Identity]*config.AlertConfig),
		Agent:     &config.AgentConfig{},
		AgentTask: make(map[config.Identity]*config.AgentTaskConfig),
		KBase:     make(map[config.Identity]*config.KnowledgeBaseConfig),
	}
	Insp = config.NewTree()
	mu   sync.RWMutex
)

func RLock() {
	mu.RLock()
}

func RUnlock() {
	mu.RUnlock()
}

func GetInsp(path config.Identity) *config.InspNode {
	mu.RLock()
	defer mu.RUnlock()
	return Insp.GetNode(path.Str())
}

func GetAllInsp() []*config.InspNode {
	mu.RLock()
	defer mu.RUnlock()
	return Insp.AllInsp
}

type ParamType interface {
	config.DefaultConfig | config.DBConfig | config.TaskConfig | config.LogConfig | config.AlertConfig |
		config.AgentConfig | config.AgentTaskConfig | config.KnowledgeBaseConfig | config.InspTree | *config.InspTree
}

type GetType interface {
	*config.DefaultConfig | *config.DBConfig | *config.TaskConfig | *config.LogConfig | *config.AlertConfig |
		*config.AgentConfig | *config.AgentTaskConfig | *config.KnowledgeBaseConfig | *config.InspTree
}

func SetConfigMeta(c config.ConfigMeta) error {
	Adds(c.CommonConfigGroup.Alerts...)
	Adds(c.CommonConfigGroup.DBs...)
	Adds(c.CommonConfigGroup.Logs...)
	Adds(c.TaskConfigGroup.Tasks...)
	Adds(c.AgentConfigGroup.Agent)
	Adds(c.AgentConfigGroup.AgentTasks...)
	Adds(c.AgentConfigGroup.KnowledgeBases...)
	return nil
}

func SetInsp(tree *config.InspTree) error {
	Insp = tree
	return nil
}

func Adds[T ParamType](configs ...T) (err error) {
	for i := range configs {
		err = Add(configs[i])
		if err != nil {
			return err
		}
	}
	return
}

func Add[T ParamType](cfg T) error {
	mu.Lock()
	defer mu.Unlock()
	switch t := any(cfg).(type) {
	case config.DBConfig:
		Meta.DBs = append(Meta.DBs, t)
		Index.DB[t.Identity] = &Meta.DBs[len(Meta.DBs)-1]
	case config.TaskConfig:
		Meta.Tasks = append(Meta.Tasks, t)
		Index.Task[t.Identity] = &Meta.Tasks[len(Meta.Tasks)-1]
	case config.LogConfig:
		Meta.Logs = append(Meta.Logs, t)
		Index.Log[t.Identity] = &Meta.Logs[len(Meta.Logs)-1]
	case config.AlertConfig:
		Meta.Alerts = append(Meta.Alerts, t)
		Index.Alert[t.Identity] = &Meta.Alerts[len(Meta.Alerts)-1]
	case config.AgentConfig:
		Meta.Agent = t
		Index.Agent = &Meta.Agent
	case *config.InspTree:
		Insp = t
	case config.AgentTaskConfig:
		Meta.AgentTasks = append(Meta.AgentTasks, t)
		Index.AgentTask[t.Identity] = &Meta.AgentTasks[len(Meta.AgentTasks)-1]
	case config.KnowledgeBaseConfig:
		Meta.KnowledgeBases = append(Meta.KnowledgeBases, t)
		Index.KBase[t.Identity] = &Meta.KnowledgeBases[len(Meta.KnowledgeBases)-1]
	default:
		return fmt.Errorf("type of config nonsupport to Add: %s\n", t)
	}
	return nil
}

func Del[T ParamType](cfg T) error {
	mu.Lock()
	defer mu.Unlock()
	switch t := any(cfg).(type) {
	//case config.DefaultConfig:
	//	Index.Default = &config.DefaultConfig{} // 非 map 类型保持清空值
	case config.DBConfig:
		delete(Index.DB, t.Identity)
		removeFromSlice[config.DBConfig](Meta.DBs, t)
	case config.TaskConfig:
		delete(Index.Task, t.Identity)
		removeFromSlice[config.TaskConfig](Meta.Tasks, t)
	case config.LogConfig:
		delete(Index.Log, t.Identity)
		removeFromSlice[config.LogConfig](Meta.Logs, t)
	case config.AlertConfig:
		delete(Index.Alert, t.Identity)
		removeFromSlice[config.AlertConfig](Meta.Alerts, t)
	case config.AgentConfig:
		Index.Agent = nil
		Meta.Agent = config.AgentConfig{}
	case config.InspTree:
		//todo：增删某个命令
		Insp = nil // 指针类型置空
	case *config.InspTree:
		Insp = nil // 指针类型置空
	case config.AgentTaskConfig:
		delete(Index.AgentTask, t.Identity)
		removeFromSlice[config.AgentTaskConfig](Meta.AgentTasks, t)
	case config.KnowledgeBaseConfig:
		delete(Index.KBase, t.Identity)
		removeFromSlice[config.KnowledgeBaseConfig](Meta.KnowledgeBases, t)
	default:
		return fmt.Errorf("type of config nonsupport to Del: %T", t) // 修正错误提示类型格式符
	}
	return nil
}

func removeFromSlice[T config.Id](slice []T, cfg T) {
	for i, item := range slice {
		if item.GetIdentity() == cfg.GetIdentity() {
			slice = append((slice)[:i], (slice)[i+1:]...)
			break
		}
	}
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
		res = any(Index.Default).(T)
	case *config.DBConfig:
		//todo: 测试是否可行
		//if db, ok := Index.DB[t.GetIdentity]; ok {
		//	res = any(db).(T)
		//} else {
		//	err = fmt.Errorf("DB config %q not found", t.GetIdentity)
		//}
		index, err := getFromIndex(Index.DB, t.Identity)
		return any(index).(T), err
	case *config.TaskConfig:
		if task, ok := Index.Task[t.Identity]; ok {
			res = any(task).(T)
		} else {
			err = fmt.Errorf("task config %q not found", t.Identity)
		}
	case *config.LogConfig:
		if log, ok := Index.Log[t.Identity]; ok {
			res = any(log).(T)
		} else {
			err = fmt.Errorf("log config %q not found", t.Identity)
		}
	case *config.AlertConfig:
		if alert, ok := Index.Alert[t.Identity]; ok {
			res = any(alert).(T)
		} else {
			err = fmt.Errorf("alert config %q not found", t.Identity)
		}
	case *config.AgentConfig:
		res = any(Index.Agent).(T)
	case *config.InspTree:
		res = any(Insp).(T) // 直接返回指针
	case *config.AgentTaskConfig:
		if task, ok := Index.AgentTask[t.Identity]; ok {
			res = any(task).(T)
		} else {
			err = fmt.Errorf("agent task %q not found", t.Identity)
		}
	case *config.KnowledgeBaseConfig:
		if kb, ok := Index.KBase[t.Identity]; ok {
			res = any(kb).(T)
		} else {
			err = fmt.Errorf("knowledge base %q not found", t.Identity)
		}
	default:
		err = fmt.Errorf("unsupported config type: %T", t)
	}
	return
}

func getFromIndex[T config.Id](index map[config.Identity]T, id config.Identity) (T, error) {
	if res, ok := index[id]; ok {
		return res, nil
	} else {
		return res, fmt.Errorf("config center: index not exist: %s in %s ", id.GetIdentity(), index)
	}
}
