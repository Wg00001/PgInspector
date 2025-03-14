package _default

import (
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
	alerter2 "PgInspector/usecase/alerter"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/18
 */

func init() {
	alerter2.RegisterDriver("default", AlerterDefault{})
}

type AlerterDefault struct {
}

func (e AlerterDefault) Init(config config.AlertConfig) (alerter.Alerter, error) {
	return AlerterDefault{}, nil
}

func (e AlerterDefault) Send(alerter.Content) error {
	return fmt.Errorf("Alert Err - Empty Alert: this alerter has not init, please check config ")
}

var _ alerter.Alerter = (*AlerterDefault)(nil)
