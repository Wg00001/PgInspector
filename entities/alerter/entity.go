package alerter

import (
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"time"
)

/**
 * @description: 监控报警功能，对insp某个具体数值进行监控，到达某个值时触发对应报警
 * @author Wg
 * @date 2025/1/19
 */

type Alerter interface {
	Send(Content) error
	Init(*config.AlertConfig) (Alerter, error)
}

type Content struct {
	TimeStamp    time.Time
	TaskName     config.Name
	DBName       config.Name
	InspName     string
	Result       db.Result
	AlertWhen    string
	AlertBecause string
}

func (c Content) AddAlertInfo(when, because string) Content {
	c.AlertWhen = when
	c.AlertBecause = because
	return c
}

func (c Content) AddWhen(when string) Content {
	c.AlertWhen = when
	return c
}

func (c Content) AddBecause(because string) Content {
	c.AlertBecause = because
	return c
}
