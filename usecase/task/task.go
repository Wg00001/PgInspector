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
		log.Println("config is nil")
		return nil
	}
	res := &task.Task{
		Config:   taskCfg,
		TargetDB: make([]*config.DBConfig, 0, len(taskCfg.TargetDB)),
		Inspects: []*insp.Node{},
	}
	for _, val := range taskCfg.TargetDB {
		res.TargetDB = append(res.TargetDB, usecase.GetDbConfig(val))
	}

	//是否全选 (全部insp)
	if taskCfg.AllInspector {
		res.Inspects = usecase.GetAllInsp()
	}
	//添加todo列表的insp
	for _, val := range taskCfg.Todo {
		res.Inspects = append(res.Inspects, usecase.GetInsp(val))
	}
	//去掉not to do的insp (使用hash连接)
	notToDo := make(map[config.Name]bool, len(taskCfg.NotTodo))
	for _, val := range taskCfg.NotTodo {
		notToDo[val] = true
	}
	newArr := make([]*insp.Node, 0, len(res.Inspects))
	for _, val := range res.Inspects {
		if !notToDo[config.Name(val.Name)] {
			newArr = append(newArr, val)
		}
	}
	res.Inspects = newArr
	return res
}
