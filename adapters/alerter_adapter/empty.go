package alerter_adapter

import (
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
)

type AlertEmpty struct {
}

func (a AlertEmpty) Send(content alerter.Content) error {
	return nil
}

func (a AlertEmpty) Init(config *config.AlertConfig) (alerter.Alerter, error) {
	return a, nil
}

var _ alerter.Alerter = (*AlertEmpty)(nil)
