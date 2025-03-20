package yaml

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	config2 "PgInspector/usecase/config"
	insp2 "PgInspector/usecase/insp"
	"PgInspector/utils"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/3/18
 */

func init() {
	config2.RegisterParser("yaml", &ConfigYamlParser{})
}

type ConfigYamlParser struct {
	ConfigYaml      ConfigYaml
	AgentConfigYaml AgentConfigYaml
	InspTree        *insp.Tree
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
	AiConfig           config.AgentConfig           `yaml:"agent"`
	AiTaskConfigOrigin []map[string]interface{}     `yaml:"agenttask"`
	AiTaskConfig       []config.AgentTaskConfig     `yaml:"-"`
	KBaseConfigOrigin  []map[string]interface{}     `yaml:"kbase"`
	KBaseConfig        []config.KnowledgeBaseConfig `yaml:"-"`
}

var _ config.Parser = (*ConfigYamlParser)(nil)

func (c *ConfigYamlParser) ParseConfig(file []byte) (_ config.CommonConfigGroup, err error) {
	file = []byte(strings.ToLower(string(file)))
	err = yaml.Unmarshal(file, &c.ConfigYaml)
	if err != nil {
		return
	}
	//处理logger设置
	c.ConfigYaml.LogConfig = make([]config.LogConfig, 0, len(c.ConfigYaml.LogConfigOrigin))
	for _, o := range c.ConfigYaml.LogConfigOrigin {
		c.ConfigYaml.LogConfig = append(c.ConfigYaml.LogConfig,
			func(origin map[string]interface{}) config.LogConfig {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("fail to read logger config")
					}
				}()
				cur := config.LogConfig{
					Identity: config.NewIdentity(origin["identity"]),
					Driver:   origin["driver"].(string),
					Header:   make(map[string]string),
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
		return
	}

	c.ConfigYaml.AlertConfig = make([]config.AlertConfig, 0, len(c.ConfigYaml.AlertConfigOrigin))
	for _, o := range c.ConfigYaml.AlertConfigOrigin {
		c.ConfigYaml.AlertConfig = append(c.ConfigYaml.AlertConfig,
			func(origin map[string]interface{}) config.AlertConfig {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("fail to read alert config")
					}
				}()
				cur := config.AlertConfig{
					Identity: config.NewIdentity(origin["identity"]),
					Driver:   origin["driver"].(string),
					Header:   make(map[string]string),
				}
				for k, v := range origin {
					if str, ok := v.(string); ok {
						cur.Header[k] = str
					}
				}
				return cur
			}(o))
	}
	return
}

func (c *ConfigYamlParser) ParseInspector(file []byte) (_ config.TaskConfigGroup, err error) {
	var origin map[string]interface{}
	err = yaml.Unmarshal(file, &origin)
	if err != nil {
		return
	}

	c.InspTree = insp.NewTree()
	var dfs func(nowPath string, node map[string]interface{}) error
	dfs = func(nowPath string, node map[string]interface{}) error {
		for k, v := range node {
			switch t := v.(type) {
			case string: //配置里直接是string说明是SQL，未配置alert
				n, err := insp2.NodeBuilder{}.WithName(k).WithSQL(t).WithEmptyAlert().Build()
				if err != nil {
					return err
				}
				err = c.InspTree.AddChild(nowPath, &n)
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
				err = c.InspTree.AddChild(nowPath, &n)
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
	return config.TaskConfigGroup{}, dfs("", origin)
}

func (c *ConfigYamlParser) ParseAgent(file []byte) (_ config.AgentConfigGroup, err error) {
	file = []byte(strings.ToLower(string(file)))
	err = yaml.Unmarshal(file, &c.AgentConfigYaml)
	if err != nil {
		return
	}

	// Process agent task configurations
	c.AgentConfigYaml.AiTaskConfig = make([]config.AgentTaskConfig, 0, len(c.AgentConfigYaml.AiTaskConfigOrigin))
	for _, origin := range c.AgentConfigYaml.AiTaskConfigOrigin {
		funcResult := func(o map[string]interface{}) (result config.AgentTaskConfig) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("fail to read agent task config: %v", r)
				}
			}()

			m := utils.UseMap(o)

			// Basic configurations
			result.Identity = config.Identity(m.GetString("identity"))
			result.LogID = config.Identity(m.GetInt("logid"))
			result.AlertID = config.Identity(m.GetInt("alertid"))
			result.KBaseResults = m.GetInt("kbaseresults")
			result.KBaseMaxLen = m.GetInt("kbasemaxlen")
			result.SystemMessage = m.GetString("systemmessage")

			// Cron configuration
			if cronMap := m.GetMap("cron"); cronMap != nil {
				result.Cron = &config.Cron{
					CronTab: cronMap.GetString("crontab"),
					Duration: func() time.Duration {
						duration, err := time.ParseDuration(cronMap.GetString("duration"))
						if err != nil {
							return 0
						}
						return duration
					}(),
					AtTime:  parseStringSlice(cronMap["attime"]),
					Weekly:  parseWeekdaySlice(cronMap["weekly"]),
					Monthly: parseIntSlice(cronMap["monthly"]),
				}
			}

			// Log filter configuration
			if logFilterMap := m.GetMap("logfilter"); logFilterMap != nil {
				result.LogFilter = config.LogFilter{
					StartTime: parseTime(logFilterMap.GetString("starttime")),
					EndTime:   parseTime(logFilterMap.GetString("endtime")),
					TaskNames: parseNames(logFilterMap["tasknames"]),
					DBNames:   parseNames(logFilterMap["dbnames"]),
					TaskIDs:   parseStringSlice(logFilterMap["taskids"]),
					InspNames: parseNames(logFilterMap["inspnames"]),
				}
			}

			// Knowledge base references
			result.KBase = parseNames(m["kbase"])
			return
		}(origin)
		c.AgentConfigYaml.AiTaskConfig = append(c.AgentConfigYaml.AiTaskConfig, funcResult)
	}
	if err != nil {
		return
	}

	// Process knowledge base configurations
	c.AgentConfigYaml.KBaseConfig = make([]config.KnowledgeBaseConfig, 0, len(c.AgentConfigYaml.KBaseConfigOrigin))
	for _, origin := range c.AgentConfigYaml.KBaseConfigOrigin {
		funcResult := func(o map[string]interface{}) (result config.KnowledgeBaseConfig) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("fail to read knowledge base config: %v", r)
				}
			}()

			m := utils.UseMap(o)
			result.Identity = config.Identity(m.GetString("name"))
			result.Driver = m.GetString("driver")
			result.Value = m
			return
		}(origin)
		c.AgentConfigYaml.KBaseConfig = append(c.AgentConfigYaml.KBaseConfig, funcResult)
	}
	return
}
