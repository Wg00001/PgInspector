package test

import (
	"PgInspector/adapters/config_reader"
	"PgInspector/entities/config"
	"PgInspector/usecase"
	"PgInspector/usecase/db"
	"PgInspector/usecase/task"
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	config.InitConfig(config_reader.NewReader("yaml", "../../app/config"))
	db.Connect(usecase.GetDbConfig("example1"))
	task.Init(usecase.GetTaskConfig("task1"))
	err := task.Do("task1")
	if err != nil {
		fmt.Println(err)
		return
	}
}
