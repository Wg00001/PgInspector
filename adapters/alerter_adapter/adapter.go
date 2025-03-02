package alerter_adapter

import (
	"PgInspector/adapters/alerter_adapter/default"
	"PgInspector/adapters/alerter_adapter/feishu"
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
		return feishu.AlerterFeishu{}.Init(config)
	default:
		return _default.AlerterDefault{}.Init(config)
	}
}
