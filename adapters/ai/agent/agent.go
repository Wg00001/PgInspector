package agent

import (
	"PgInspector/entities/ai"
	"PgInspector/entities/config"
	ai2 "PgInspector/usecase/ai"
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"log"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/5
 */

func init() {
	ai2.RegisterDriver("agent", AnalyzerAgent{})
}

type AnalyzerAgent config.AiConfig

var _ ai.Analyzer = (*AnalyzerAgent)(nil)

func (a AnalyzerAgent) Init(aiConfig *config.AiConfig) (ai.Analyzer, error) {
	return AnalyzerAgent(*aiConfig), nil
}

func (a AnalyzerAgent) Analyze(s string) (string, error) {
	llm, err := openai.New(openai.WithBaseURL(a.Api), openai.WithToken(a.ApiKey), openai.WithModel(a.Model))
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	msg := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem,
			`你是一个拥有20年经验的DBA专家，严格按以下JSON格式要求输出分析结果：
				{
				  "巡检指标数据分析": [
					{
					  "数据库指标": "指标名称",
					  "指标状态": "正常/需调整",
					  "指标变化趋势": "上升/下降/稳定",
					  "潜在问题": "简明问题描述",
					  "性能优化建议": "具体操作建议"
					}
				  ],
				  "补充建议": [
					"扩展性建议1",
					"架构优化建议2"
				  ]
				}
				
				请确保：
				1. 数值型指标必须包含单位（如128MB）
				2. 状态判断使用【正常】【需调整】两种分类
				3. 变化趋势需基于时间序列数据判断`),
		llms.TextParts(llms.ChatMessageTypeHuman, "请分析以下日志:\n"+s),
	}
	log.Printf("发送消息: %v\n", msg)
	resp, err := llm.GenerateContent(ctx, msg)
	if err != nil {
		return "", err
	}
	log.Printf("收到消息：%v\n", resp.Choices[0].Content)
	return resp.Choices[0].Content, nil
}
