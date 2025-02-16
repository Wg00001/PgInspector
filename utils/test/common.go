package test

import (
	"PgInspector/adapters/config_adapter"
	"PgInspector/adapters/logger_adapter"
	"PgInspector/entities/config"
	"PgInspector/usecase"
	"PgInspector/usecase/db"
	"PgInspector/usecase/logger"
	"PgInspector/usecase/task"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

func initConfig() {
	config.InitConfig(config_adapter.NewReader("yaml", "../../app/config"))

}

func initDB() {
	d, err := db.InitDB(usecase.GetDbConfig("example1"))
	if err != nil {
		panic(err)
	}
	err = db.Register(d)
	if err != nil {
		panic(err)
	}
}

func initTask() {
	t1, err := task.InitTask(usecase.GetTaskConfig("task1"))
	if err != nil {
		panic(err)
	}
	err = task.Register(t1)
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	lg, err := logger_adapter.NewLogger(usecase.GetLoggerConfig(1))
	if err != nil {
		panic(err)
	}
	err = logger.Register(lg)
	if err != nil {
		panic(err)
	}
}
