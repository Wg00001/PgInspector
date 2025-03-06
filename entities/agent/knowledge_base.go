package agent

import "PgInspector/entities/config"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/5
 */

type Document struct {
	ID        string
	Content   string
	Embedding []float32
	Metadata  map[string]interface{}
}

type KnowledgeBase interface {
	Init(config *config.KnowledgeBaseConfig) (KnowledgeBase, error)
	WriteIn(docs []*Document) error
	Search(topK int, query ...string) ([]*Document, error)
	SimilaritySearch(topK int, embedding []float32) ([]*Document, error)
	Embedding(query string) ([]float32, error)
}
