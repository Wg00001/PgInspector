package db

import (
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"fmt"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/10
 */

var pool = sync.Map{}

func Register(sqlDB *db.SqlDB) error {
	if _, ok := pool.Load(sqlDB.Config.Name); ok {
		return fmt.Errorf("sql db is already exist, db name: %s\n", sqlDB.Config.Name)
	}
	pool.Store(sqlDB.Config.Name, sqlDB)
	return nil
}

func Get(arg config.Name) *db.SqlDB {
	if val, ok := pool.Load(arg); ok {
		return val.(*db.SqlDB)
	}
	return &db.SqlDB{Err: fmt.Errorf("db config is nil")}
}

func Close(arg config.Name) error {
	if val, ok := pool.LoadAndDelete(arg); ok {
		err := val.(*db.SqlDB).Close()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("DB not exist")
	}
	return nil
}

func CloseAll() error {
	var err error
	pool.Range(func(key, value any) bool {
		er := Close(key.(config.Name))
		if er != nil {
			err = er
		}
		return true
	})
	return err
}
