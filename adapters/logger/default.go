package logger

import (
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

func (d DefaultLogger) Gout(rows *sql.Rows) {
	utils.PrintQuery(rows)
}

var _ logger.Logger = (*DefaultLogger)(nil)
