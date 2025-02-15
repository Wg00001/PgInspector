package logger

import (
	"PgInspector/entities/config"
	"PgInspector/entities/logger"
	"PgInspector/utils"
	"database/sql"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type DefaultLogger struct {
}

func (d DefaultLogger) GetID() config.ID {
	return 0
}

func (d DefaultLogger) Log(l logger.InspLog, rows *sql.Rows) {
	utils.PrintQuery(l, rows)
}

var _ logger.Logger = (*DefaultLogger)(nil)
