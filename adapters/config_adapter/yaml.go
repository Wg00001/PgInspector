package config_adapter

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"PgInspector/usecase"
	"PgInspector/utils"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

/**
 * @description: insp
 * @author Wg
 * @date 2025/1/19
 */

type ConfigReaderYaml struct {
	FilePath   string
	ConfigName string
	InspName   string
	cyaml      ConfigYaml
	insp       *insp.Tree
}

type ConfigYaml struct {
	Default         config.DefaultConfig     `yaml:"default"`
	DBConfigs       []config.DBConfig        `yaml:"db"`
	TaskConfigs     []config.TaskConfig      `yaml:"task"`
	LogConfigOrigin []map[string]interface{} `yaml:"log"`
	LogConfig       []config.LogConfig       `yaml:"-"`
	AlertConfig     []config.AlertConfig     `yaml:"alert"`
}

var _ config.Reader = (*ConfigReaderYaml)(nil)

func (c *ConfigReaderYaml) ReadConfig() error {
	file, err := os.ReadFile(c.FilePath + c.ConfigName)
	if err != nil {
		return err
	}
	file = []byte(strings.ToLower(string(file)))
	err = yaml.Unmarshal(file, &c.cyaml)
	if err != nil {
		return err
	}

	//处理logger设置
	c.cyaml.LogConfig = make([]config.LogConfig, 0, len(c.cyaml.LogConfigOrigin))
	for _, o := range c.cyaml.LogConfigOrigin {
		c.cyaml.LogConfig = append(c.cyaml.LogConfig,
			func(origin map[string]interface{}) config.LogConfig {
				defer func() {
					if r := recover(); r != nil {
						log.Println("fail to read logger config")
					}
				}()
				cur := config.LogConfig{
					LogID:  config.ID(origin["id"].(int)),
					Driver: origin["driver"].(string),
					Header: make(map[string]string),
				}
				for k, v := range origin {
					if str, ok := v.(string); ok {
						cur.Header[k] = str
					}
				}
				return cur
			}(o))
	}
	return nil
}

func (c *ConfigReaderYaml) ReadInspector() error {
	file, err := os.ReadFile(c.FilePath + c.InspName)
	if err != nil {
		return err
	}
	var origin map[string]interface{}
	err = yaml.Unmarshal(file, &origin)
	if err != nil {
		return err
	}

	c.insp = insp.NewTree()
	var dfs func(nowPath string, node map[string]interface{})
	dfs = func(nowPath string, node map[string]interface{}) {
		for k, v := range node {
			switch t := v.(type) {
			case map[string]interface{}:
				c.insp.AddChild(nowPath, &insp.Node{Name: k})
				dfs(nowPath+k, t)
			case string:
				c.insp.AddChild(nowPath, &insp.Node{Name: k, SQL: t})
			}
		}
	}
	dfs("", origin)
	return nil
}

func (c *ConfigReaderYaml) SaveIntoConfig() {
	usecase.AddConfigs(c.cyaml.DBConfigs...)
	usecase.AddConfigs(c.cyaml.TaskConfigs...)
	usecase.AddConfigs(c.cyaml.LogConfig...)
	usecase.AddConfigs(c.cyaml.AlertConfig...)
	usecase.AddConfigs(c.insp)
}

func (c *ConfigReaderYaml) FormatFilename() {
	if strings.Index(c.FilePath, "/") != len(c.FilePath)-1 {
		c.FilePath += "/"
	}
	if c.ConfigName == "" {
		c.ConfigName = "config.yaml"
	}
	if c.InspName == "" {
		c.InspName = "inspect.yaml"
	}
	c.ConfigName = utils.FileNameFormat(c.ConfigName, "yaml")
	c.InspName = utils.FileNameFormat(c.InspName, "yaml")
}
