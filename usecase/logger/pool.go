package logger

import (
	"PgInspector/entities/config"
	"PgInspector/entities/logger"
	"fmt"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

var (
	pool = make(map[config.ID]logger.Logger)
	mu   sync.Mutex
)

func Register(lg logger.Logger) error {
	if pool == nil {
		return fmt.Errorf("init logger_adapter fail: logger_adapter pool is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	pool[lg.GetID()] = lg
	return nil
}

func Get(id config.ID) logger.Logger {
	mu.Lock()
	defer mu.Unlock()
	if val, ok := pool[id]; ok {
		return val
	}
	return nil
}
