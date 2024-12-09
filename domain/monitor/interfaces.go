package monitor

//将监控等抽象成接口,便于对接新的数据库

// Monitor 每个monitor对象对应一个数据或数据表
type Monitor interface {
	InitConfig()
}

// DBMonitor 数据库监控的抽象接口
type DBMonitor interface {
	// QPS 和 TPS
	MonitorDBQueryRate() float64       // 监控查询 QPS
	MonitorDBTransactionRate() float64 // 监控事务 TPS
	// 线程状态
	MonitorDBActiveConnections() int // 监控活动连接数
	MonitorDBMaxConnections() int    // 监控最大连接数
	MonitorDBThreadCount() int       // 监控线程数
	// 慢查询和锁等待
	MonitorDBSlowQueryCount() int // 监控慢查询次数
	MonitorDBLockWaitCount() int  // 监控等待锁的次数
	// 磁盘 I/O
	MonitorDBReadIO() int64  // 监控磁盘读字节数
	MonitorDBWriteIO() int64 // 监控磁盘写字节数
	MonitorDBIOOps() int64   // 监控磁盘 I/O 操作数
}

// TableMonitor 表监控的抽象接口
type TableMonitor interface {
	// 数据表状态
	MonitorTableSize(tableName string) float64             // 监控表数据大小
	MonitorTableGrowth(tableName string) float64           // 监控表的增长率
	MonitorTableAutoIncrementStatus(tableName string) bool // 监控表主键自增状态
}
