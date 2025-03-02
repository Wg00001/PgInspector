package alerter

import (
	"PgInspector/entities/alerter"
	"fmt"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/2
 */

var (
	drivers = make(map[string]alerter.Alerter)
	mu      sync.RWMutex
)

func RegisterDriver(name string, driver alerter.Alerter) {
	mu.Lock()
	defer mu.Unlock()
	if drivers == nil {
		panic("alerter: drivers map is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("alerter: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func GetDriver(name string) (alerter.Alerter, error) {
	mu.RLock()
	defer mu.RUnlock()
	res, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("alerter: get driver fail %s\n", name)
	}
	return res, nil
}
