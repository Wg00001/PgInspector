package monitor

import (
	"PgInspector/pkg/config"
	"database/sql"
	"fmt"
)

type BaseMonitor struct {
	DB *sql.DB
	*config.DbConfig
	Error error
}

var _ Monitor = (*BaseMonitor)(nil)

func (m *BaseMonitor) InitConfig() {
	if m.DbConfig.DSN != "" {
		m.Error = fmt.Errorf("DSN is empty, but no driver implement to get default value")
		return
	}
	if m.DbConfig.DriverName == "" {
		m.Error = fmt.Errorf("DriverName is empty,but no driver implement to get default value")
		return
	}
	m.DB, m.Error = sql.Open(m.DbConfig.DriverName, m.DbConfig.DSN)
	if m.Error != nil {
		return
	}
	m.Error = m.DB.Ping()
}
