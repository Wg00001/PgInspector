package postgre

import (
	"PgInspector/pkg"
)

// 实现interfaces中关于监控数据库函数
type DBMonitorImpl struct {
	pkg.MonitorImpl
}

var _ pkg.DBMonitor = (*DBMonitorImpl)(nil)

func (m DBMonitorImpl) MonitorDBQueryRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBTransactionRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBActiveConnections() int {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBMaxConnections() int {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBThreadCount() int {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBSlowQueryCount() int {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBLockWaitCount() int {
	//TODO implement me
	panic("implement me")
}

func (m DBMonitorImpl) MonitorDBReadIO() int64 {
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
