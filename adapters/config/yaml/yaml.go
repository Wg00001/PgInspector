package yaml

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	config2 "PgInspector/usecase/config"
	insp2 "PgInspector/usecase/insp"
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
	cyaml      ConfigYaml
	acyaml     AgentConfigYaml
	insp       *insp.Tree
}

type ConfigYaml struct {
	Default           config.DefaultConfig     `yaml:"default"`
	DBConfigs         []config.DBConfig        `yaml:"db"`
	TaskConfigs       []config.TaskConfig      `yaml:"task"`
	LogConfigOrigin   []map[string]interface{} `yaml:"log"`
	LogConfig         []config.LogConfig       `yaml:"-"`
	AlertConfigOrigin []map[string]interface{} `yaml:"alert"`
	AlertConfig       []config.AlertConfig     `yaml:"-"`
}

type AgentConfigYaml struct {
	AiConfig     config.AgentConfig           `yaml:"agent"`
	AiTaskConfig []config.AgentTaskConfig     `yaml:"agenttask"`
	KBaseConfig  []config.KnowledgeBaseConfig `yaml:"knowledgebase"`
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
	}, nil
}

func (c *ConfigReaderYaml) ReadConfig() (err error) {
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
						err = fmt.Errorf("fail to read logger config")
					}
				}()
				cur := config.LogConfig{
					ID:     config.ID(origin["id"].(int)),
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
	if err != nil {
		return err
	}

	c.cyaml.AlertConfig = make([]config.AlertConfig, 0, len(c.cyaml.AlertConfigOrigin))
	for _, o := range c.cyaml.AlertConfigOrigin {
		c.cyaml.AlertConfig = append(c.cyaml.AlertConfig,
			func(origin map[string]interface{}) config.AlertConfig {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("fail to read alert config")
					}
				}()
				cur := config.AlertConfig{
					ID:     config.ID(origin["id"].(int)),
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
	if err != nil {
		return err
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
	var dfs func(nowPath string, node map[string]interface{}) error
	dfs = func(nowPath string, node map[string]interface{}) error {
		for k, v := range node {
			switch t := v.(type) {
			case string: //配置里直接是string说明是SQL，未配置alert
				n, err := insp2.NodeBuilder{}.WithName(k).WithSQL(t).WithEmptyAlert().Build()
				if err != nil {
					return err
				}
				err = c.insp.AddChild(nowPath, &n)
				if err != nil {
					return err
				}

			case map[string]interface{}: //是map说明有配置alert, 进行解析
				nb, err := ParseMap(insp2.NodeBuilder{}.WithName(k), t)
				if err != nil {
					return err
				}
				n, err := nb.Build()
				if err != nil {
					return err
				}
				err = c.insp.AddChild(nowPath, &n)
				if err != nil {
					return err
				}
				err = dfs(nowPath+k, t)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
	return dfs("", origin)
}

func (c *ConfigReaderYaml) ReadAgent() error {
	file, err := os.ReadFile(c.FilePath + c.AgentName)
	if err != nil {
		return err
	}
	file = []byte(strings.ToLower(string(file)))
	err = yaml.Unmarshal(file, &c.acyaml)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigReaderYaml) SaveIntoConfig() {
	config2.AddConfigs(c.cyaml.DBConfigs...)
	config2.AddConfigs(c.cyaml.TaskConfigs...)
	config2.AddConfigs(c.cyaml.LogConfig...)
	config2.AddConfigs(c.cyaml.AlertConfig...)
	config2.AddConfigs(c.acyaml.AiConfig)
	config2.AddConfigs(c.acyaml.AiTaskConfig...)
	config2.AddConfigs(c.acyaml.KBaseConfig...)
	config2.AddConfigs(c.insp)
}

func ParseMap(n insp2.NodeBuilder, arg map[string]interface{}) (m insp2.NodeBuilder, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("inspect node build fail, please check inspect config \nerr: %v\n", r)
		}
	}()
	alertId, ok := arg[keyAlertId]
	if !ok {
		return n, nil
	} else {
		n.AlertID = config.ID(alertId.(int))
		delete(arg, keyAlertId)
	}
	alertWhen, ok := arg[keyAlertWhen]
	if !ok {
		return n, nil
	} else {
		delete(arg, keyAlertWhen)
	}
	n = n.BuildAlertFunc(alertWhen.(string))

	sql, ok := arg[keySQL]
	if !ok {
		return n, nil
	} else {
		n.SQL = sql.(string)
		delete(arg, keySQL)
	}
	return n, nil
}
