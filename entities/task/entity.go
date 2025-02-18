package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"PgInspector/entities/logger"
	"PgInspector/usecase/db"
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
	logger.Logger
}

func (t *Task) Do() error {
	for _, inspect := range t.Inspects {
		for _, tdb := range t.TargetDB {
			if tdb == nil {
				continue
			}
			query, err := db.Get(tdb.Name).Query(inspect.SQL)
			if err != nil {
				return err
			}
			result, err := insp.RowsToMap(query)
			if err != nil {
				return err
			}
			err = inspect.AlertFunc(result)
			if err != nil {
				return err
			}
			t.Logger.Log(logger.InspLog{
				Timestamp: time.Now(),
				TaskName:  t.Config.GetName(),
				DBName:    tdb.Name,
			}, result)
		}
	}
	log.Printf("task finish: %s\n", t.Config.TaskName)
	return nil
}
