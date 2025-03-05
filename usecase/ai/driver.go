package ai

import (
	"PgInspector/entities/ai"
	"PgInspector/entities/config"
	"fmt"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/5
 */

func Use(aiConfig config.AiConfig) error {
	driver, err := GetDriver(aiConfig.Driver)
	if err != nil {
		return err
	}
	init, err := driver.Init(&aiConfig)
	if err != nil {
		return err
	}
	Register(init)
	return nil
}

var (
	drivers  = make(map[string]ai.Analyzer)
	muDriver sync.RWMutex
)

func RegisterDriver(name string, driver ai.Analyzer) {
	muDriver.Lock()
	defer muDriver.Unlock()
	if drivers == nil {
		panic("ai: drivers map is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("ai: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func GetDriver(name string) (ai.Analyzer, error) {
	muDriver.RLock()
	defer muDriver.RUnlock()
	res, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("ai: get driver fail %s\n", name)
	}
	return res, nil
}
