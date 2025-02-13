package logger

import "database/sql"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type Logger interface {
	Gout(*sql.Rows)
}
