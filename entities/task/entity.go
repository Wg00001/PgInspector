package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"PgInspector/entities/logger"
	"PgInspector/usecase/db"
	"PgInspector/utils"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type Task struct {
	Config   *config.TaskConfig
	TargetDB []*config.DBConfig
	Inspects []*insp.Node
	logger.Logger
}

func (t *Task) Do() error {
	logFunc := utils.PrintQuery
	if t.Logger != nil {
		logFunc = t.Logger.Gout
	}
	for _, inspect := range t.Inspects {
		for _, tdb := range t.TargetDB {
			if tdb == nil {
				continue
			}
			query, err := db.Get(tdb.Name).Query(inspect.SQL)
			if err != nil {
				return err
			}
			logFunc(query)
		}
	}
	return nil
}
