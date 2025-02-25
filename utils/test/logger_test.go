package test

import (
	"PgInspector/adapters/cron"
	"PgInspector/adapters/start"
	"PgInspector/usecase/task"
	"testing"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

func TestLogger(t *testing.T) {
	initConfig()
	initDB("example1")
	initDB("example2")
	initLogger()
	initTask()
	cron.Init()
	cron.AddTask(task.Get("task1"))
	cron.Start()
	select {}
}

func TestStart(t *testing.T) {
	start.SetConfigPath("../../app/config", "yaml")
	start.Init()
	start.Run()
}
