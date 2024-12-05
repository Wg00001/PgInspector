package task

type DBMonitorData struct {
	DBName            string
	DSN               string  //连接该数据库的DSN,账号密码需要进行加密
	QueryRate         float64 // 查询速率
	TransactionRate   float64 // 事务速率
	ActiveConnections int     // 活跃连接数
	MaxConnections    int     // 最大连接数
	ThreadCount       int     // 数据库线程数量
	SlowQueryCount    int     // 慢查询数
	LockWaitCount     int     // 锁等待数
	ReadIO            int64   // 读 IO 数量
	WriteIO           int64   // 写 IO 数量
	IOOps             int64   // IO 操作总数
}
