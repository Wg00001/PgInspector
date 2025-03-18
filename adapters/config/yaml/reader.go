package yaml

import (
	"PgInspector/entities/config"
	config2 "PgInspector/usecase/config"
	"fmt"
	"os"
	"strings"
)

/**
 * @description: insp
 * @author Wg
 * @date 2025/1/19
 */

const (
	optionFilepath   = "filepath"
	optionConfigName = "config"
	optionInspName   = "inspect"
	optionAgentName  = "agent"

	keyAlertId   = "_alertId"
	keyAlertWhen = "_alertWhen"
	keySQL       = "_sql"
)

func init() {
	config2.RegisterDriver("yaml", &ConfigReaderYaml{})
}

type ConfigReaderYaml struct {
	FilePath   string
	ConfigName string
	InspName   string
	AgentName  string
	parser     *ConfigYamlParser
}

var _ config.Reader = (*ConfigReaderYaml)(nil)

func (c *ConfigReaderYaml) NewReader(option map[string]string) (_ config.Reader, err error) {
	filepath, ok := option[optionFilepath]
	if !ok {
		return nil, fmt.Errorf("config reader: option deficiency - %s\n", optionFilepath)
	}
	configName, ok := option[optionConfigName]
	if !ok {
		configName = optionConfigName + ".yaml"
	}
	inspName, ok := option[optionInspName]
	if !ok {
		inspName = optionInspName + ".yaml"
	}
	agentName, ok := option[optionAgentName]
	if !ok {
		agentName = optionAgentName + ".yaml"
	}

	if !strings.HasSuffix(filepath, "/") {
		filepath += "/"
	}
	return &ConfigReaderYaml{
		FilePath:   filepath,
		ConfigName: configName,
		InspName:   inspName,
		AgentName:  agentName,
		parser:     &ConfigYamlParser{},
	}, nil
}

func (c *ConfigReaderYaml) ReadConfig() (err error) {
	file, err := os.ReadFile(c.FilePath + c.ConfigName)
	if err != nil {
		return err
	}
	return c.parser.ParseConfig(file)
}

func (c *ConfigReaderYaml) ReadInspector() error {
	file, err := os.ReadFile(c.FilePath + c.InspName)
	if err != nil {
		return err
	}
	return c.parser.ParseInsp(file)
}

func (c *ConfigReaderYaml) ReadAgent() error {
	file, err := os.ReadFile(c.FilePath + c.AgentName)
	if err != nil {
		return err
	}
	return c.parser.ParseAgent(file)
}
func (c *ConfigReaderYaml) SaveIntoConfig() {
	config2.Sets(c.parser.ConfigYaml.DBConfigs...)
	config2.Sets(c.parser.ConfigYaml.TaskConfigs...)
	config2.Sets(c.parser.ConfigYaml.LogConfig...)
	config2.Sets(c.parser.ConfigYaml.AlertConfig...)
	config2.Sets(c.parser.AgentConfigYaml.AiConfig)
	config2.Sets(c.parser.AgentConfigYaml.AiTaskConfig...)
	config2.Sets(c.parser.AgentConfigYaml.KBaseConfig...)
	config2.Sets(c.parser.InspTree)
}

func (c *ConfigReaderYaml) Watch() {
	//TODO implement me
	panic("implement me")
}
