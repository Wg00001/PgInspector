package ai

import (
	"PgInspector/entities/ai"
	"PgInspector/entities/config"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

//根据driver找到对应的adapter实现，以init全局analyzer

var a ai.Analyzer

func Init(cfg *config.AiConfig) {
	t, err := a.Init(cfg)
	if err != nil {
		return
	}
	a = t
}

func Analyze(input string) (string, error) {
	if a == nil {
		return "", fmt.Errorf("ai analyzer has not init")
	}
	return a.Analyze(input)
}
