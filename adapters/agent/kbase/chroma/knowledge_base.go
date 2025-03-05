package chroma

import (
	"PgInspector/entities/agent"
	"PgInspector/entities/config"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/5
 */

type KBaseChroma struct {
	Config     config.KnowledgeBaseConfig
	Path       string //path是chroma的文件存储位置，存储于本地
	Collection string //chroma的collection类似于库
}

var _ agent.KnowledgeBase = (*KBaseChroma)(nil)

func (k KBaseChroma) Init(config *config.KnowledgeBaseConfig) (_ agent.KnowledgeBase, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("kbase: chroma init fail - panic: %s\n", r)
		}
	}()
	k.Config = *config
	k.Path = config.Value["path"]
	k.Collection = config.Value["collection"]
	return k, nil
}

func (k KBaseChroma) WriteIn(docs []agent.Document) error {
	//TODO implement me
	panic("implement me")
}

func (k KBaseChroma) Search(query string, topK int) ([]agent.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (k KBaseChroma) SimilaritySearch(embedding []float32, topK int) ([]agent.Document, error) {
	//TODO implement me
	panic("implement me")
}
