package ollama

import (
	"PgInspector/entities/agent"
	"PgInspector/entities/config"
	ai2 "PgInspector/usecase/agent"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
	"log"
	"strings"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

func init() {
	ai2.RegisterDriver("ollama", AnalyzerOllama{})
}

type AnalyzerOllama config.AgentConfig

func (a AnalyzerOllama) Init(aiConfig *config.AgentConfig) (agent.Analyzer, error) {
	return AnalyzerOllama(*aiConfig), nil
}

func (a AnalyzerOllama) Analyze(s string) (string, error) {
	opt := llm.SetOptions(map[string]interface{}{
		option.Temperature: a.Temperature,
	})
	question := llm.GenQuery{
		Model:   a.Model,
		Prompt:  "这是我的数据库巡检日志，请进行分析，并对数据库运行状态给出简短的评价和建议：\n" + s,
		Options: opt,
	}
	log.Printf("发送日志：%v\n", question)
	answer, err := completion.Generate(a.Api, question)
	return withoutThink(answer.Response), err
}

var _ agent.Analyzer = (*AnalyzerOllama)(nil)

func withoutThink(s string) string {
	if strings.Contains(s, "</think>") {
		return strings.Split(s, "</think>")[1]
	}
	return s
}
