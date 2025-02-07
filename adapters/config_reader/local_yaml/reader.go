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
 * @description: insp
 * @author Wg
 * @date 2025/1/19
 */

type ConfigReaderYaml struct {
	FilePath string
	cyaml    ConfigYaml
}

type ConfigYaml struct {
	Default       config.DefaultConfig `yaml:"default"`
	DBConfigs     []config.DBConfig    `yaml:"db"`
	TaskConfigs   []config.TaskConfig  `yaml:"task"`
	LogConfig     []config.LogConfig   `yaml:"log"`
	AlertConfig   []config.AlertConfig `yaml:"alert"`
	InspectConfig map[string]string    `yaml:"inspect"`
	//InspectConfig interface{} `yaml:"inspect"`
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
	usecase.AddConfigs(c.cyaml.LogConfig...)
	usecase.AddConfigs(c.cyaml.AlertConfig...)
	insp := make([]config.InspectConfig, 0, len(c.cyaml.InspectConfig))
	for name, sql := range c.cyaml.InspectConfig {
		insp = append(insp, config.InspectConfig{
			InspName: config.Name(name),
			SQL:      sql,
		})
	}
	usecase.AddConfigs(insp...)
	fmt.Printf("%+v\n", c.cyaml)
	fmt.Printf("%+v\n", usecase.GetConfig())
}
