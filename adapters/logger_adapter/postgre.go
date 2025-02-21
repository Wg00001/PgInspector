package logger_adapter

import (
	"PgInspector/entities/config"
	db2 "PgInspector/entities/db"
	"PgInspector/entities/logger"
	"PgInspector/usecase/db"
	"database/sql"
	"fmt"
	"log"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type LogPostgre struct {
	Config       *config.LogConfig
	LogDBName    config.Name
	LogTableName config.Name
}

func (l LogPostgre) Init(cfg *config.LogConfig) (logger.Logger, error) {
	dbName, ok := cfg.Header["dbname"]
	if !ok {
		return LogPostgre{}, fmt.Errorf("Log target db is not exist, dbName:%s\n", dbName)
	}
	tableName, ok := cfg.Header["tablename"]
	if !ok {
		tableName = "inspect_log"
	}
	return LogPostgre{Config: cfg, LogDBName: config.Name(dbName), LogTableName: config.Name(tableName)}, nil

}

func (l LogPostgre) GetID() config.ID {
	return l.Config.LogID
}

func (l LogPostgre) Log(inspLog logger.Content, result db2.Result) {
	//res, err := insp.RowsToMap(rows)
	//if err != nil {
	//	l.Output(inspLog.WithErr(err))
	//	return
	//}
	l.Output(inspLog.WithJSON(result))
}

func (l LogPostgre) Output(res logger.Content) {
	db.Get(l.LogDBName)
	// 获取数据库连接
	logDB := db.Get(l.LogDBName)
	if logDB.Err != nil {
		log.Printf("Failed to get database connection: %v", logDB.Err)
		return
	}

	if l.LogTableName == "" {
		l.LogTableName = "inspect_log"
		// 检查表是否存在
		exists, err := checkTableExists(logDB.DB, l.LogTableName)
		if err != nil {
			log.Printf("Failed to check table existence: %v", err)
			return
		}

		// 如果表不存在，则创建表
		if !exists {
			err = createTable(logDB.DB, l.LogTableName)
			if err != nil {
				log.Printf("Failed to create table: %v", err)
				return
			}
		}
	}

	// 构建插入 SQL 语句
	insertQuery := fmt.Sprintf(`
        INSERT INTO %s (timestamp, task_name, db_name, task_id, result)
        VALUES ($1, $2, $3, $4, $5)
    `, l.LogTableName)

	// 执行插入操作
	_, err := logDB.DB.Exec(insertQuery, res.Timestamp, res.TaskName, res.DBName, res.TaskID, res.Result)
	if err != nil {
		log.Printf("Failed to insert log data: %v", err)
	}
}

var _ logger.Logger = (*LogPostgre)(nil)

// checkTableExists 检查指定表是否存在
func checkTableExists(db *sql.DB, tableName config.Name) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE  table_schema = 'public'
        AND    table_name   = $1
    )`
	err := db.QueryRow(query, tableName.Str()).Scan(&exists)
	return exists, err
}

// createTable 创建指定表
func createTable(db *sql.DB, tableName config.Name) error {
	createQuery := fmt.Sprintf(`
        CREATE TABLE %s (
            id SERIAL PRIMARY KEY,
            timestamp TIMESTAMP,
            task_name TEXT,
            db_name TEXT,
            task_id TEXT,
            result JSONB
        )
    `, tableName.Str())
	_, err := db.Exec(createQuery)
	return err
}
