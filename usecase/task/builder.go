package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"PgInspector/entities/task"
	"PgInspector/usecase"
	"PgInspector/usecase/logger"
	"fmt"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/10
 */

func InitTask(taskCfg *config.TaskConfig) (res *task.Task, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("init task fail: %s", err.Error())
		}
	}()
	if taskCfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	log := logger.Get(taskCfg.LogLevel)
	if log == nil {
		return nil, fmt.Errorf("taskname:  %s\ntask init err: log id not exist, please check log config or init log before init task\n", taskCfg.TaskName)
	}
	res = &task.Task{
		Identity: taskCfg.TaskName.Str() + time.Now().Format(time.RFC3339),
		Config:   taskCfg,
		TargetDB: make([]*config.DBConfig, 0, len(taskCfg.TargetDB)),
		Inspects: []*insp.Node{},
		Logger:   log,
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
	return res, nil
}
