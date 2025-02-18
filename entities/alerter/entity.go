package alerter

import (
	"PgInspector/entities/config"
)

/**
 * @description: 监控报警功能，对insp某个具体数值进行监控，到达某个值时触发对应报警
 * @author Wg
 * @date 2025/1/19
 */

type Alert struct {
	*config.AlertConfig
	Alerter //insp发现超过阈值后，调用此接口来发送alert
}

type Alerter interface {
	Send(any)
}
