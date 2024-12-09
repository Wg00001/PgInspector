package monitor

import (
	"PgInspector/domain/config"
	"PgInspector/domain/monitor"
	"database/sql"
	"fmt"
)

// 实现interfaces中关于监控数据库函数
type PostgreMonitor struct {
	monitor.BaseMonitor
}

var _ monitor.Monitor = (*PostgreMonitor)(nil)
var _ monitor.DBMonitor = (*PostgreMonitor)(nil)

func (m *PostgreMonitor) InitConfig() {
	if m.DBConfig.DSN != "" {
		m.Error = fmt.Errorf("DSN is empty, but no driver implement to get default value")
		return
	}
	if m.DBConfig.Driver == "" {
		m.DBConfig.Driver = config.DiverNamePostgreSQL
	}
	m.DB, m.Error = sql.Open(m.DBConfig.Driver, m.DBConfig.DSN)
	if m.Error != nil {
		return
	}
	m.Error = m.DB.Ping()
}

func (m *PostgreMonitor) MonitorDBQueryRate() float64 {
	return monitor.Query[float64](&m.BaseMonitor, `
		SELECT coalesce(sum(xact_commit + xact_rollback)::float / extract(epoch from now() - pg_postmaster_start_time()), 0)
		FROM pg_stat_database;
	`)
}
func (m *PostgreMonitor) MonitorDBTransactionRate() float64 {
	return monitor.Query[float64](&m.BaseMonitor, `
		SELECT coalesce(sum(xact_commit + xact_rollback)::float / extract(epoch from now() - pg_postmaster_start_time()), 0)
		FROM pg_stat_database;
	`)
}

func (m *PostgreMonitor) MonitorDBActiveConnections() int {
	return monitor.Query[int](&m.BaseMonitor, `
		SELECT count(*)
		FROM pg_stat_activity
		WHERE state = 'active';
	`)
}

func (m *PostgreMonitor) MonitorDBMaxConnections() int {
	return monitor.Query[int](&m.BaseMonitor, `
		SELECT setting::int
		FROM pg_settings
		WHERE name = 'max_connections';
	`)
}

func (m *PostgreMonitor) MonitorDBThreadCount() int {
	return monitor.Query[int](&m.BaseMonitor, `
		SELECT count(*)
		FROM pg_stat_activity;
	`)
}

func (m *PostgreMonitor) MonitorDBSlowQueryCount() int {
	return monitor.Query[int](&m.BaseMonitor, `
		SELECT count(*)
		FROM pg_stat_activity
		WHERE now() - query_start > interval '1 minute' AND state = 'active';
	`)
}

func (m *PostgreMonitor) MonitorDBLockWaitCount() int {
	return monitor.Query[int](&m.BaseMonitor, `
		SELECT count(*)
		FROM pg_locks
		WHERE granted = false;
	`)
}

func (m *PostgreMonitor) MonitorDBReadIO() int64 {
	return monitor.Query[int64](&m.BaseMonitor, `
		SELECT sum(blks_read)
		FROM pg_stat_database;
	`)
}

func (m *PostgreMonitor) MonitorDBWriteIO() int64 {
	return monitor.Query[int64](&m.BaseMonitor, `
		SELECT sum(blks_written)
		FROM pg_stat_database;
	`)
}

func (m *PostgreMonitor) MonitorDBIOOps() int64 {
	return monitor.Query[int64](&m.BaseMonitor, `
		SELECT sum(blks_read + blks_written)
		FROM pg_stat_database;
	`)
}
