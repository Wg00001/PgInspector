package etcd

import (
	"PgInspector/adapters/config/parser/yaml"
	"PgInspector/entities/config"
	config2 "PgInspector/usecase/config"
	"context"
	"github.com/coreos/etcd/clientv3"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/18
 */

type ConfigReaderEtcd struct {
	client      *clientv3.Client
	rootContext context.Context
	cfg         config.Config
	parser      yaml.ConfigYamlParser
}

var _ config.Reader = (*ConfigReaderEtcd)(nil)

func (c ConfigReaderEtcd) NewReader(option map[string]string) (_ config.Reader, err error) {
	c.client, err = clientv3.New(clientv3.Config{})
	return c, err
}

func (c ConfigReaderEtcd) ReadConfig() error {
	_, err := c.client.Get(c.rootContext, "config")
	if err != nil {
		return err
	}
	return nil
}

func (c ConfigReaderEtcd) ReadInspector() error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigReaderEtcd) ReadAgent() error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigReaderEtcd) SaveIntoConfig() {
	config2.Sets(c.parser.ConfigYaml.DBConfigs...)
	config2.Sets(c.parser.ConfigYaml.TaskConfigs...)
	config2.Sets(c.parser.ConfigYaml.LogConfig...)
	config2.Sets(c.parser.ConfigYaml.AlertConfig...)
	config2.Sets(c.parser.AgentConfigYaml.AiConfig)
	config2.Sets(c.parser.AgentConfigYaml.AiTaskConfig...)
	config2.Sets(c.parser.AgentConfigYaml.KBaseConfig...)
	config2.Sets(c.parser.InspTree)
}

func (c ConfigReaderEtcd) Watch() {
	//TODO implement me
	panic("implement me")
}
