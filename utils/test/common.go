package test

import (
	"PgInspector/adapters/config_reader"
	"PgInspector/entities/config"
	"PgInspector/usecase"
	"PgInspector/usecase/db"
	"PgInspector/usecase/task"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

// init config,db,task
func initTest() {
	config.InitConfig(config_reader.NewReader("yaml", "../../app/config"))

	initDB, err := db.InitDB(usecase.GetDbConfig("example1"))
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Register(initDB)

	t1, err := task.InitTask(usecase.GetTaskConfig("task1"))
	if err != nil {
		fmt.Println(err)
		return
	}
	task.Register(t1)
}
