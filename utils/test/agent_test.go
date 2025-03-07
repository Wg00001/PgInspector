package test

import (
	_ "PgInspector/adapters/agent/kbase/chroma"
	"PgInspector/adapters/start"
	"PgInspector/entities/agent"
	"PgInspector/entities/config"
	"PgInspector/usecase/agent/kbase"
	"fmt"
	"strings"
	"testing"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/6
 */

func TestKBase(t *testing.T) {
	start.SetConfigPath("../../app/config", "yaml")
	start.Init()

	// 初始化测试配置
	cfg := config.KnowledgeBaseConfig{
		Name:   "unit-test-kb",
		Driver: "chroma",
		Value: map[string]string{
			"path":       "http://localhost:8000",
			"collection": "default",
			"tenant":     "default",
		},
	}

	// 测试初始化
	err := kbase.Use(cfg)
	if err != nil {
		t.Fatal(err)
	}
	kb := kbase.Get("unit-test-kb")
	fmt.Printf("%+v\n", kb)

	// 准备测试数据
	testDocs := []*agent.Document{
		{
			ID:        "doc1",
			Content:   "Go语言并发编程指南",
			Embedding: []float32{0.1, 0.2, 0.3},
			Metadata:  map[string]interface{}{"author": "张三"},
		},
		{
			ID:        "doc2",
			Content:   "分布式系统设计原则",
			Embedding: []float32{0.4, 0.5, 0.6},
			Metadata:  map[string]interface{}{"year": 2023},
		},
	}

	t.Run("WriteDocuments", func(t *testing.T) {
		// 正常写入测试
		if err := kb.WriteIn(testDocs); err != nil {
			t.Fatalf("WriteIn failed: %v", err)
		}

		// 写入空数据测试
		t.Run("EmptyDocuments", func(t *testing.T) {
			err := kb.WriteIn(nil)
			if err == nil {
				t.Error("Expected error for empty documents")
			}
		})

		// 无效文档测试
		t.Run("InvalidDocument", func(t *testing.T) {
			invalidDocs := []*agent.Document{{ID: ""}}
			err := kb.WriteIn(invalidDocs)
			if err == nil {
				t.Error("Expected error for invalid document")
			}
		})
	})

	t.Run("Search", func(t *testing.T) {
		// 正常搜索测试
		t.Run("BasicSearch", func(t *testing.T) {
			results, err := kb.Search(2, "并发")
			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}

			if len(results) != 1 {
				t.Errorf("Expected 1 result, got %d", len(results))
			}

			if !strings.Contains(results[0].Content, "并发") {
				t.Errorf("Unexpected result content: %s", results[0].Content)
			}
		})

		// 多查询测试
		t.Run("MultipleQueries", func(t *testing.T) {
			results, err := kb.Search(2, "系统", "设计")
			if err != nil {
				t.Fatal(err)
			}

			if len(results) < 1 {
				t.Error("No results found for multiple queries")
			}
		})

		// 边界测试
		t.Run("BoundaryConditions", func(t *testing.T) {
			// 零结果测试
			results, err := kb.Search(0, "无关内容")
			if err != nil || len(results) != 0 {
				t.Error("TopK=0 should return empty results")
			}

			// 无效查询测试
			_, err = kb.Search(1, "")
			if err == nil {
				t.Error("Expected error for empty query")
			}
		})
	})

	t.Run("SimilaritySearch", func(t *testing.T) {
		queryEmbedding := []float32{0.12, 0.22, 0.32} // 接近doc1的向量

		// 正常相似度搜索
		results, err := kb.SimilaritySearch(1, queryEmbedding)
		if err != nil {
			t.Fatal(err)
		}

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].ID != "doc1" {
			t.Errorf("Expected doc1, got %s", results[0].ID)
		}

		// 错误输入测试
		t.Run("InvalidInput", func(t *testing.T) {
			// 空向量测试
			_, err := kb.SimilaritySearch(1, nil)
			if err == nil {
				t.Error("Expected error for empty embedding")
			}

			// 维度不匹配测试（假设维度为3）
			_, err = kb.SimilaritySearch(1, []float32{0.1})
			if err == nil {
				t.Error("Expected error for dimension mismatch")
			}
		})
	})

	t.Run("Embedding", func(t *testing.T) {
		// 正常生成测试
		emb, err := kb.Embedding("生成向量")
		if err != nil {
			t.Fatal(err)
		}

		// 验证向量维度
		if len(emb) != 3 { // 根据测试数据维度
			t.Errorf("Expected 3-dim embedding, got %d", len(emb))
		}

		// 空输入测试
		_, err = kb.Embedding("")
		if err == nil {
			t.Error("Expected error for empty input")
		}
	})

	// 清理测试数据
	t.Cleanup(func() {
		if memKb, ok := kb.(interface{ Reset() }); ok {
			memKb.Reset()
		}
	})
}
