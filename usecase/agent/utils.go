package agent

import (
	"PgInspector/entities/agent"
	alerter2 "PgInspector/entities/alerter"
	"PgInspector/entities/config"
	"PgInspector/usecase/agent/analyzer"
	"encoding/json"
	"fmt"
	"github.com/wg00001/wgo-sdk/wg"
	"strings"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/6
 */

func buildAiAlertContent(t *AgentTask, msg string) *alerter2.Content {
	return &alerter2.Content{
		TimeStamp: time.Now(),
		TaskName:  "AiAnalyzeTask: " + t.Name,
		TaskID:    time.Now().Format("20060504_150201"),
		InspName:  logFilterString(t.LogFilter),
		Result:    []map[string]interface{}{{"message": msg}},
	}
}

func logFilterString(filter config.LogFilter) string {
	marshal, err := json.Marshal(filter)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func generateQueryEmbedding(query string, base agent.KnowledgeBase) ([]float32, error) {
	// 生成嵌入向量（假设有嵌入生成器）
	embedding, err := base.Embedding(query)
	if err != nil {
		return nil, fmt.Errorf("embedding generation failed: %w", err)
	}
	return embedding, nil
}

// 使用Ai进行分析获取搜索关键词
func generateQueryWithAI(logContent *string) (string, error) {
	res, err := analyzer.Analyze(&agent.AnalyzeContent{
		SystemMsg: "",
		UserMsg:   "请从以下日志中提取3到5个最相关的搜索关键词，并用逗号分隔的关键词列表。\n日志内容：\n" + *logContent,
		KBaseMsg:  "",
	})
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(res, ",", " "), nil
}

// 混合检索 - 结合关键词和向量检索
func hybridSearch(kbase agent.KnowledgeBase, query string, embedding []float32, topK int) ([]*agent.Document, error) {
	// 并行执行两种检索
	keywordResults, _ := kbase.Search(topK/2, query)
	vectorResults, _ := kbase.SimilaritySearch(topK/2, embedding)
	// 结果去重合并
	idx := wg.SliceToSet(keywordResults,
		func(document *agent.Document) string {
			return document.ID
		})
	for _, vr := range vectorResults {
		if _, ok := idx[vr.ID]; !ok {
			keywordResults = append(keywordResults, vr)
		}
	}
	return keywordResults, nil
}

// 生成知识上下文（辅助函数）
func formatKBaseContent(docs []*agent.Document, maxLen int) *string {
	var buf strings.Builder
	currentLen := 0

	for _, doc := range docs {
		content := fmt.Sprintf("条目 %s:\n%s\n\n", doc.ID, doc.Content)
		if currentLen+len(content) > maxLen {
			break
		}
		buf.WriteString(content)
		currentLen += len(content)
	}

	if buf.Len() == 0 {
		return nil
	}
	s := buf.String()
	return &s
}
