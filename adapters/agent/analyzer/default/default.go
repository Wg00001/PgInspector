package _default

import (
	"PgInspector/entities/agent"
	"PgInspector/entities/config"
	ai2 "PgInspector/usecase/agent/analyzer"
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

func (a AnalyzerDefault) Init(config *config.AgentConfig) (agent.Analyzer, error) {
	//TODO implement me
	panic("implement me")
}

func (a AnalyzerDefault) Analyze(s *agent.AnalyzeContent) (string, error) {
	//TODO implement me
	panic("implement me")
}

var _ agent.Analyzer = (*AnalyzerDefault)(nil)
