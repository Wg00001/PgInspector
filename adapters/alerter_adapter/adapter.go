package alerter_adapter

import (
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/18
 */

func NewAlerter(config *config.AlertConfig) (alerter.Alerter, error) {
	switch config.Driver {
	case "feishu":
		return AlerterFeishu{}.Init(config)
	default:
		return AlerterDefault{}.Init(config)
	}
}
