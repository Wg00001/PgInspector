package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	config2 "PgInspector/usecase/config"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/10
 */

// NewTask
// alert如果没有设置，那么应该继承父节点的alertID
func NewTask(taskCfg *config.TaskConfig) (res *Task, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("init task fail: %s", err.Error())
		}
	}()
	if taskCfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	res = &Task{
		//Identity: taskCfg.Name.Str() + time.Now().Format(time.RFC3339),
		Config:   taskCfg,
		TargetDB: make([]*config.DBConfig, 0, len(taskCfg.TargetDB)),
		Inspects: []*insp.Node{},
	}
	for _, val := range taskCfg.TargetDB {
		dbcfg, err := config2.Get[*config.DBConfig](&config.DBConfig{Name: val})
		if err != nil {
			return nil, err
		}
		res.TargetDB = append(res.TargetDB, dbcfg)
	}

	//是否全选 (全部insp)
	if taskCfg.AllInspector {
		res.Inspects = config2.GetAllInsp()
	}
	//添加todo列表的insp
	for _, val := range taskCfg.Todo {
		res.Inspects = append(res.Inspects, config2.GetInsp(val))
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
