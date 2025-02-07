package db

import (
	"PgInspector/entities/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/6
 */

type SqlDB struct {
	Config *config.DBConfig
	DB     *sql.DB
	err    error
}

func (c *SqlDB) Connect() {
	c.DB, c.err = sql.Open(c.Config.Driver, c.Config.DSN)
}

func (c *SqlDB) Error() error {
	return c.err
}

func Connect(dbConfig *config.DBConfig) (*SqlDB, error) {
	if dbConfig == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	cur := &SqlDB{Config: dbConfig}
	cur.Connect()
	return cur, cur.err
}
