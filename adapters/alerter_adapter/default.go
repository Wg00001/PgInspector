package alerter_adapter

import (
	"PgInspector/entities/alerter"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/18
 */

type EmptyAlert struct {
}

func (e EmptyAlert) Send(alerter.Content) error {
	return fmt.Errorf("Alert Err - Empty Alert: this alerter has not init, please check config ")
}

var _ alerter.Alerter = (*EmptyAlert)(nil)
