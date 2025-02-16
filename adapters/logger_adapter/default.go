package logger_adapter

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

type LogDefault struct {
}

func (d LogDefault) GetID() config.ID {
	return 0
}

func (d LogDefault) Log(l logger.InspLog, rows *sql.Rows) {
	utils.PrintQuery(l, rows)
}

var _ logger.Logger = (*LogDefault)(nil)
