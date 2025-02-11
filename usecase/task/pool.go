package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/task"
	"PgInspector/usecase"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/11
 */

var pool = sync.Map{}

func Init(arg *config.TaskConfig) *task.Task {
	cur := BuildTask(usecase.GetTaskConfig(config.GetNameT(arg)))
	pool.Store(cur.Config.TaskName, cur)
	return cur
}

func Get(name config.Name) *task.Task {
	if val, ok := pool.Load(name); ok {
		return val.(*task.Task)
	}
	return nil
}

func Do() {

}

func Delete(name config.Name) {
	pool.Delete(name)
}
