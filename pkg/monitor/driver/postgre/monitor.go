package postgre

import (
	"PgInspector/pkg/config"
	"PgInspector/pkg/monitor"
	"database/sql"
	"fmt"
)

// 实现interfaces中关于监控数据库函数
type DBMonitorImpl struct {
	monitor.BaseMonitor
}

var _ monitor.Monitor = (*DBMonitorImpl)(nil)
var _ monitor.DBMonitor = (*DBMonitorImpl)(nil)

func (m *DBMonitorImpl) InitConfig() {
	if m.DbConfig.DSN != "" {
		m.DbConfig.DSN = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			m.DbConfig.Host, m.DbConfig.Port, m.DbConfig.Username, m.DbConfig.Password, m.DbConfig.Dbname)
	}
	if m.DbConfig.DriverName == "" {
		m.DbConfig.DriverName = config.DiverNamePostgreSQL
	}
	m.DB, m.Error = sql.Open(m.DbConfig.DriverName, m.DbConfig.DSN)
	if m.Error != nil {
		return
	}
	m.Error = m.DB.Ping()
}

func (m *DBMonitorImpl) MonitorDBQueryRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBTransactionRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBActiveConnections() int {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBMaxConnections() int {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBThreadCount() int {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBSlowQueryCount() int {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBLockWaitCount() int {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBReadIO() int64 {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBWriteIO() int64 {
	//TODO implement me
	panic("implement me")
}

func (m *DBMonitorImpl) MonitorDBIOOps() int64 {
	//TODO implement me
	panic("implement me")
}

//todo:初始化方式
//todo:实现接口
