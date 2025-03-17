package task

import (
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
	db2 "PgInspector/entities/db"
	"PgInspector/entities/insp"
	"PgInspector/entities/logger"
	"PgInspector/entities/task"
	"PgInspector/usecase/db"
	logger2 "PgInspector/usecase/logger"
	"context"
	"fmt"
	"log"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type Task struct {
	//Id string //批次编号, task每次启动会生成一个
	Config   *config.TaskConfig
	TargetDB []*config.DBConfig
	Inspects []*insp.Node
}

var _ task.Task = (*Task)(nil)

func (t *Task) Do(ctx context.Context) error {
	taskid := time.Now().Format("20060102_150405")
	fmt.Printf("task: start - %s\n", taskid)
	for _, inspect := range t.Inspects {
		for _, tdb := range t.TargetDB {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			if tdb == nil {
				continue
			}
			//执行SQL
			query, err := db.Get(tdb.Identity).Query(inspect.SQL)
			if err != nil {
				return err
			}
			result, err := db2.RowsToResult(query)
			if err != nil {
				return err
			}

			//记录
			logger2.Get(t.Config.LogID).Log(logger.Content{
				Timestamp: time.Now(),
				TaskName:  t.Config.Identity,
				TaskID:    taskid,
				InspName:  inspect.Name,
				DBName:    tdb.Identity,
				Result:    result,
			})

			//报警
			err = inspect.AlertFunc(alerter.Content{
				TimeStamp: time.Now(),
				TaskName:  t.Config.Identity,
				TaskID:    taskid,
				DBName:    tdb.Identity,
				InspName:  inspect.Name,
				Result:    result,
			})
			if err != nil {
				return err
			}

		}
	}
	log.Printf("task finish: %s\n", t.Config.Identity)
	return nil
}

func (t *Task) GetCron() *config.Cron {
	return t.Config.Cron
}

func (t *Task) Identity() config.Identity {
	return "insp_task:" + t.Config.Identity
}
