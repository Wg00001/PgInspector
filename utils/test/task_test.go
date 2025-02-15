package test

import (
	"PgInspector/adapters/task/cron"
	"PgInspector/usecase/task"
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	initTest()
	fmt.Println(task.Do("task1"))
}

func TestCron(t *testing.T) {
	initTest()
	cron.Init()
	cron.AddTask(task.Get("task1"))
	cron.Start()
	for {
	}
}
