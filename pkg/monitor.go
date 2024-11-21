package pkg

import (
	"PgInspector/pkg/config"
	"database/sql"
	"fmt"
)

type MonitorImpl struct {
	DB *sql.DB
	*config.DbConfig
	Error error
}

func (m MonitorImpl) InitConfig() {
	m.Error = fmt.Errorf("monitor without driver")
}

var _ Monitor = (*MonitorImpl)(nil)
