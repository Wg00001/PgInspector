package logger_adapter

import (
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"PgInspector/entities/logger"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type LogDefault struct {
}

func (d LogDefault) Init(cfg *config.LogConfig) (logger.Logger, error) {
	return LogDefault{}, nil
}

func (d LogDefault) GetID() config.ID {
	return 0
}

func (d LogDefault) Log(l logger.Content, result db.Result) {
	//utils.PrintQuery(l, rows)
	fmt.Println(result)
}

var _ logger.Logger = (*LogDefault)(nil)
