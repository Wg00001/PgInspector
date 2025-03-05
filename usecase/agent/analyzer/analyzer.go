package analyzer

import (
	"PgInspector/entities/agent"
	"fmt"
	"log"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

//根据driver找到对应的adapter实现，以init全局analyzer

var (
	a  agent.Analyzer
	mu sync.Mutex
)

func Register(oa agent.Analyzer) {
	mu.Lock()
	defer mu.Unlock()
	a = oa
	log.Printf("openai: registry: %#v\n", oa)
}

func Analyze(input string) (string, error) {
	if a == nil {
		return "", fmt.Errorf("openai analyzer has not init")
	}
	return a.Analyze(input)
}
