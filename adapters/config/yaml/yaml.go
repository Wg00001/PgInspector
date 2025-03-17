package yaml

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	config2 "PgInspector/usecase/config"
	insp2 "PgInspector/usecase/insp"
	"PgInspector/utils"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
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
	AiConfig           config.AgentConfig           `yaml:"agent"`
	AiTaskConfigOrigin []map[string]interface{}     `yaml:"agenttask"`
	AiTaskConfig       []config.AgentTaskConfig     `yaml:"-"`
	KBaseConfigOrigin  []map[string]interface{}     `yaml:"kbase"`
	KBaseConfig        []config.KnowledgeBaseConfig `yaml:"-"`
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

	// Process agent task configurations
	c.acyaml.AiTaskConfig = make([]config.AgentTaskConfig, 0, len(c.acyaml.AiTaskConfigOrigin))
	for _, origin := range c.acyaml.AiTaskConfigOrigin {
		funcResult := func(o map[string]interface{}) (result config.AgentTaskConfig) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("fail to read agent task config: %v", r)
				}
			}()

			m := utils.UseMap(o)

			// Basic configurations
			result.Name = config.Name(m.GetString("name"))
			result.LogID = config.ID(m.GetInt("logid"))
			result.AlertID = config.ID(m.GetInt("alertid"))
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
		c.acyaml.AiTaskConfig = append(c.acyaml.AiTaskConfig, funcResult)
	}
	if err != nil {
		return err
	}

	// Process knowledge base configurations
	c.acyaml.KBaseConfig = make([]config.KnowledgeBaseConfig, 0, len(c.acyaml.KBaseConfigOrigin))
	for _, origin := range c.acyaml.KBaseConfigOrigin {
		funcResult := func(o map[string]interface{}) (result config.KnowledgeBaseConfig) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("fail to read knowledge base config: %v", r)
				}
			}()

			m := utils.UseMap(o)
			result.Name = config.Name(m.GetString("name"))
			result.Driver = m.GetString("driver")
			result.Value = m
			return
		}(origin)
		c.acyaml.KBaseConfig = append(c.acyaml.KBaseConfig, funcResult)
	}
	return err
}
func (c *ConfigReaderYaml) SaveIntoConfig() {
	config2.Sets(c.cyaml.DBConfigs...)
	config2.Sets(c.cyaml.TaskConfigs...)
	config2.Sets(c.cyaml.LogConfig...)
	config2.Sets(c.cyaml.AlertConfig...)
	config2.Sets(c.acyaml.AiConfig)
	config2.Sets(c.acyaml.AiTaskConfig...)
	config2.Sets(c.acyaml.KBaseConfig...)
	config2.Sets(c.insp)
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

// Helper functions with panic for type conversion errors
func parseTime(s interface{}) time.Time {
	str, ok := s.(string)
	if !ok || s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		t, err = time.Parse(time.TimeOnly, str)
	}
	if err != nil {
		panic(err)
	}
	return t
}

func parseNames(s interface{}) []config.Name {
	items, ok := s.([]interface{})
	if !ok {
		return nil
	}
	names := make([]config.Name, 0, len(items))
	for _, item := range items {
		names = append(names, config.Name(item.(string)))
	}
	return names
}

func parseStringSlice(data interface{}) []string {
	if data == nil {
		return nil
	}
	items, ok := data.([]interface{})
	if !ok {
		return nil
	}
	strs := make([]string, 0, len(items))
	for _, item := range items {
		strs = append(strs, fmt.Sprintf("%v", item))
	}
	return strs
}

func parseIntSlice(data interface{}) []int {
	if data == nil {
		return nil
	}
	items, ok := data.([]interface{})
	if !ok {
		return nil
	}
	ints := make([]int, 0, len(items))
	for _, item := range items {
		ints = append(ints, item.(int))
	}
	return ints
}

func parseWeekdaySlice(data interface{}) []time.Weekday {
	ints := parseIntSlice(data)
	res := make([]time.Weekday, len(ints))
	for i, v := range ints {
		res[i] = time.Weekday(v)
	}
	return res
}
