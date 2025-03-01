package logger_adapter

import (
	"PgInspector/entities/config"
	"PgInspector/entities/logger"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

const (
	Default    = "default"
	PostgreSQL = "postgres"
)

func NewLogger(cfg *config.LogConfig) (logger.Logger, error) {
	if cfg == nil {
		return nil, fmt.Errorf("init logger err, logger config is nil")
	}
	return NewLoggerWithDriver(cfg.Driver, cfg)
}

func NewLoggerWithDriver(driver string, cfg *config.LogConfig) (logger.Logger, error) {
	switch driver {
	case Default, "":
		return LogDefault{}, nil
	case PostgreSQL:
		return LogPostgre{}.Init(cfg)
	}
	return nil, fmt.Errorf("get logger_adapter err, driver is not exist: %s\n", driver)
}
