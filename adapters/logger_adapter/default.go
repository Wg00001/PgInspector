package logger_adapter

import (
	"PgInspector/entities/config"
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

func (d LogDefault) ReadLog(filter logger.Filter) ([]logger.Content, error) {
	return nil, fmt.Errorf("default logger can't read, pease use other driver")
}

func (d LogDefault) Init(cfg *config.LogConfig) (logger.Logger, error) {
	return LogDefault{}, nil
}

func (d LogDefault) GetID() config.ID {
	return 0
}

func (d LogDefault) Log(l logger.Content) {
	//utils.PrintQuery(l, rows)
	fmt.Println(l.Result)
}

var _ logger.Logger = (*LogDefault)(nil)
