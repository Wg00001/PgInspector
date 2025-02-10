package db

import (
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"PgInspector/usecase"
	"fmt"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/10
 */

var pool = sync.Map{}

func Connect[T config.Name | config.DBConfig | string](arg T) *db.SqlDB {
	dbConfig := usecase.GetDbConfig(config.GetNameT(arg))
	if dbConfig == nil {
		return &db.SqlDB{Err: fmt.Errorf("db config is nil")}
	}
	cur := &db.SqlDB{Config: dbConfig}
	cur.Connect()
	pool.Store(dbConfig.Name, cur)
	return cur
}

func Get[T config.Name | config.DBConfig | string](arg T) *db.SqlDB {
	if val, ok := pool.Load(config.GetNameT(arg)); ok {
		return val.(*db.SqlDB)
	}
	return nil
}

func Delete[T config.Name | config.DBConfig | string](arg T) {
	pool.Delete(config.GetNameT(arg))
}
