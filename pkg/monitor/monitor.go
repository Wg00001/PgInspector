package monitor

import (
	"PgInspector/pkg/config"
	"database/sql"
	"fmt"
	"github.com/wg00001/wgo-sdk/wg_log"
)

type BaseMonitor struct {
	DB *sql.DB
	*config.DBConfig
	Error error
}

var _ Monitor = (*BaseMonitor)(nil)

func (m *BaseMonitor) InitConfig() {
	if m.DBConfig.DSN != "" {
		m.Error = fmt.Errorf("DSN is empty, but no driver implement to get default value")
	}
	if m.DBConfig.Driver == "" {
		m.Error = fmt.Errorf("DriverName is empty,but no driver implement to get default value")
	}
	wg_log.FatalIfErr(m.Error)
	m.DB, m.Error = sql.Open(m.DBConfig.Driver, m.DBConfig.DSN)
	if m.Error != nil {
		return
	}
	m.Error = m.DB.Ping()
}

func Query[T float64 | int | int8 | int64](bm *BaseMonitor, query string, args ...string) T {
	var res T
	err := bm.DB.QueryRow(query, args).Scan(&res)
	if err != nil {
		bm.Error = err
	}
	return res
}
