package db

import (
	"PgInspector/entities/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/6
 */

type SqlDB struct {
	*sql.DB
	Config *config.DBConfig
	Err    error
}

func (c *SqlDB) Connect() {
	c.DB, c.Err = sql.Open(c.Config.Driver, c.Config.DSN)
}

func (c *SqlDB) Error() error {
	return c.Err
}
