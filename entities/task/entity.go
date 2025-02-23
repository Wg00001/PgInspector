package task

import (
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
	db2 "PgInspector/entities/db"
	"PgInspector/entities/insp"
	"PgInspector/entities/logger"
	"PgInspector/usecase/db"
	logger2 "PgInspector/usecase/logger"
	"log"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type Task struct {
	Identity string //批次编号, task每次启动会生成一个
	Config   *config.TaskConfig
	TargetDB []*config.DBConfig
	Inspects []*insp.Node
}

func (t *Task) Do() error {
	taskid := time.Now().Format("20060102_150405")
	for _, inspect := range t.Inspects {
		for _, tdb := range t.TargetDB {
			if tdb == nil {
				continue
			}
			//执行SQL
			query, err := db.Get(tdb.Name).Query(inspect.SQL)
			if err != nil {
				return err
			}
			result, err := db2.RowsToMap(query)
			if err != nil {
				return err
			}

			//记录
			logger2.Get(t.Config.LogID).Log(logger.Content{
				Timestamp: time.Now(),
				TaskName:  t.Config.GetName(),
				TaskID:    taskid,
				InspName:  inspect.Name,
				DBName:    tdb.Name,
			}, result)

			//报警
			err = inspect.AlertFunc(alerter.Content{
				TimeStamp: time.Now(),
				TaskName:  t.Config.GetName(),
				TaskID:    taskid,
				DBName:    tdb.Name,
				InspName:  inspect.Name,
				Result:    result,
			})
			if err != nil {
				return err
			}
		}
	}
	log.Printf("task finish: %s\n", t.Config.TaskName)
	return nil
}
