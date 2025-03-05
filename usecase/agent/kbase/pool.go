package kbase

import (
	"PgInspector/entities/agent"
	"PgInspector/entities/config"
	"fmt"
	"sync"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/5
 */

var pool = sync.Map{}

func Register(name config.Name, lg agent.KnowledgeBase) error {
	if _, ok := pool.Load(name); ok {
		return fmt.Errorf("agent kbase register fail: kbase is already exsit, name repeat - %s\n", name)
	}
	pool.Store(name, lg)
	return nil
}

func Get(name config.Name) agent.KnowledgeBase {
	val, ok := pool.Load(name)
	if !ok {
		res, _ := GetDriver("default")
		return res
	}
	t, ok := val.(agent.KnowledgeBase)
	if !ok {
		res, _ := GetDriver("default")
		return res
	}
	return t
}
