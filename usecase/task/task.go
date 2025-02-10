package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"PgInspector/entities/task"
	"PgInspector/usecase"
	"log"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/10
 */

func BuildTask(taskCfg *config.TaskConfig) *task.Task {
	if taskCfg == nil {
		log.Println("config is null")
		return nil
	}
	res := &task.Task{
		Config:   taskCfg,
		TargetDB: make([]*config.DBConfig, 0, len(taskCfg.TargetDB)),
		Inspects: []*insp.Inspect{},
	}
	for _, val := range taskCfg.TargetDB {
		res.TargetDB = append(res.TargetDB, usecase.GetDbConfig(val))
	}
	return res
}
