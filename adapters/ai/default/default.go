package _default

import (
	"PgInspector/entities/ai"
	"PgInspector/entities/config"
	ai2 "PgInspector/usecase/ai"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/1
 */

func init() {
	ai2.RegisterDriver("default", AnalyzerDefault{})
}

type AnalyzerDefault struct {
}

func (a AnalyzerDefault) Init(config *config.AiConfig) (ai.Analyzer, error) {
	//TODO implement me
	panic("implement me")
}

func (a AnalyzerDefault) Analyze(s string) (string, error) {
	//TODO implement me
	panic("implement me")
}

var _ ai.Analyzer = (*AnalyzerDefault)(nil)
