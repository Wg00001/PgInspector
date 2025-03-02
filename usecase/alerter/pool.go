package alerter

import (
	"PgInspector/adapters/alerter_adapter/default"
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

func Use(alertConfig config.AlertConfig) error {
	driver, err := GetDriver(alertConfig.Driver)
	if err != nil {
		return err
	}
	init, err := driver.Init(alertConfig)
	if err != nil {
		return err
	}
	return Register(alertConfig.AlertID, init)
}

var pool = sync.Map{}

func Register(id config.ID, alert alerter.Alerter) error {
	if _, ok := pool.Load(id); ok {
		return fmt.Errorf("alerter registry fail: alert is already exist, alert id repeat")
	}
	pool.Store(id, alert)
	return nil
}

func GetAlert(id config.ID) alerter.Alerter {
	val, ok := pool.Load(id)
	if !ok {
		return _default.AlerterDefault{}
	}
	t, ok := val.(alerter.Alerter)
	if !ok {
		return _default.AlerterDefault{}
	}
	return t
}
