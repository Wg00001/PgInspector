package local_yaml

import (
	"PgInspector/entities/config"
	"PgInspector/usecase"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

/**
 * @description: inspector
 * @author Wg
 * @date 2025/1/19
 */

type ConfigReaderYaml struct {
	FilePath string
	cyaml    ConfigYaml
}

type ConfigYaml struct {
	Default      config.DefaultConfig `yaml:"defaultconfig"`
	DBConfigs    []config.DBConfig    `yaml:"dbconfig"`
	TableConfigs []config.TableConfig `yaml:"tableconfig"`
	TaskConfigs  []config.TaskConfig  `yaml:"taskconfig"`
	LogConfig    []config.LogConfig   `yaml:"logconfig"`
	AlertConfig  []config.AlertConfig `yaml:"alertconfig"`
}

var _ config.Reader = (*ConfigReaderYaml)(nil)

func (c *ConfigReaderYaml) ReadFromSource() error {
	file, err := os.ReadFile(c.FilePath)
	if err != nil {
		return err
	}
	file = []byte(strings.ToLower(string(file)))
	err = yaml.Unmarshal(file, &c.cyaml)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigReaderYaml) SaveIntoConfig() {
	usecase.InitConfig()
	usecase.AddConfigs(c.cyaml.DBConfigs...)
	usecase.AddConfigs(c.cyaml.TaskConfigs...)
	usecase.AddConfigs(c.cyaml.TableConfigs...)
	usecase.AddConfigs(c.cyaml.LogConfig...)
	usecase.AddConfigs(c.cyaml.AlertConfig...)
	fmt.Printf("%+v\n", c.cyaml)
	fmt.Printf("%+v", usecase.GetConfig())
}
