package test

import (
	"PgInspector/adapters/task/cron"
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
	initDB()
	initTask()
	cron.Init()
	cron.AddTask(task.Get("task1"))
	cron.Start()
	for {
	}
}
