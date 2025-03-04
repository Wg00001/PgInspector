package db

import (
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

func InitDB(dbConfig *config.DBConfig) (*db.SqlDB, error) {
	if dbConfig == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	cur := &db.SqlDB{Config: dbConfig}
	cur.Connect()
	if err := cur.Ping(); err != nil {
		return nil, err
	} else {
		fmt.Println("    db: connected - " + dbConfig.Name)
	}
	return cur, nil
}
