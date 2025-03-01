package ai

import (
	"PgInspector/entities/ai"
	"PgInspector/entities/config"
	"strings"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

func NewAiAnalyzer(config *config.AiConfig) (ai.Analyzer, error) {
	switch strings.ToLower(config.Driver) {
	case "ollama":
		return AnalyzerOllama{}.Init(config)
	default:
		return AnalyzerDefault{}.Init(config)
	}
}
