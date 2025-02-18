package alerter

import (
	"PgInspector/adapters/alerter_adapter"
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
	"fmt"
	"sync"
)

/**
 * @description: 读取配置后，alert在此处注册，insp执行后根据对应的id和目标来检查
 * @author Wg
 * @date 2025/2/15
 */

var pool = sync.Map{}

func Registry(id config.ID, alert alerter.Alerter) error {
	if _, ok := pool.Load(id); ok {
		return fmt.Errorf("alerter registry fail: alert is already exist, alert id repeat")
	}
	pool.Store(id, alert)
	return nil
}

func GetAlert(id config.ID) alerter.Alerter {
	val, ok := pool.Load(id)
	if !ok {
		return alerter_adapter.EmptyAlert{}
	}
	t, ok := val.(alerter.Alerter)
	if !ok {
		return alerter_adapter.EmptyAlert{}
	}
	return t
}
