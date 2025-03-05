package agent

import "PgInspector/entities/config"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/5
 */

type Document struct {
	ID       string
	Content  string
	Metadata map[string]interface{}
}

type KnowledgeBase interface {
	Init(config *config.KnowledgeBaseConfig) (KnowledgeBase, error)
	WriteIn(docs []Document) error
	Search(query string, topK int) ([]Document, error)
	SimilaritySearch(embedding []float32, topK int) ([]Document, error)
}
