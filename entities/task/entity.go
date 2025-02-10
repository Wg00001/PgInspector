package task

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type Task struct {
	Config   *config.TaskConfig
	TargetDB []*config.DBConfig
	Inspects []*insp.Inspect
}
