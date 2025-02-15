package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/task"
	"fmt"
	"sync"
)

/**
 * @description: task pool
 * @author Wg
 * @date 2025/2/11
 */

var pool = sync.Map{}

func Register(t *task.Task) error {
	if t == nil {
		return fmt.Errorf("task init err, build task fail")
	}
	pool.Store(t.Config.TaskName, t)
	return nil
}

func Get(name config.Name) *task.Task {
	if val, ok := pool.Load(name); ok {
		return val.(*task.Task)
	}
	return nil
}

func Delete(name config.Name) {
	pool.Delete(name)
}

func Do(name config.Name) error {
	t := Get(name)
	if t == nil {
		return fmt.Errorf("name of task not exist, name: %s\n", name)
	}
	return t.Do()
}
