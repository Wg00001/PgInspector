package config

import (
	"PgInspector/entities/config"
	"fmt"
	"log"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/4
 */

func Open(driverName string, option map[string]string) error {
	//todo: option封装
	reader, err := GetDriver(driverName)
	if err != nil {
		return err
	}
	reader, err = reader.NewReader(option)
	if err != nil {
		return err
	}
	err = reader.ReadConfig()
	if err != nil {
		return err
	}
	err = reader.ReadInspector()
	if err != nil {
		return err
	}
	reader.SaveIntoConfig()
	log.Println("config initiated...")
	return nil
}

var (
	drivers  = make(map[string]config.Reader)
	driverMu sync.RWMutex
)

func RegisterDriver(name string, driver config.Reader) {
	driverMu.Lock()
	defer driverMu.Unlock()
	if drivers == nil {
		panic("config: drivers map is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("config: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func GetDriver(name string) (config.Reader, error) {
	driverMu.RLock()
	defer driverMu.RUnlock()
	res, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("config: get driver fail %s\n", name)
	}
	return res, nil
}
