package test

import (
	"PgInspector/adapters/ai"
	"PgInspector/adapters/cron"
	"PgInspector/adapters/start"
	"PgInspector/entities/config"
	ai2 "PgInspector/usecase/ai"
	config2 "PgInspector/usecase/config"
	"fmt"
	"testing"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/26
 */

func TestOllamaApi(t *testing.T) {
	a := ai.AnalyzerOllama{
		Driver:      "ollama",
		Api:         "http://127.0.0.1:11434",
		ApiKey:      "",
		Model:       "deepseek-r1:7b",
		Temperature: 0.5,
	}
	analyze, err := a.Analyze("你好，你是什么模型")
	if err != nil {
		return
	}
	fmt.Println(analyze)
}

func TestAiTask(t *testing.T) {
	//start.SetConfigPath("../../app/config", "yaml")
	fmt.Println(config2.Open("yaml", map[string]string{
		"filepath": "../../app/config",
	}))
	fmt.Println(start.InitDB())
	cron.Init()
	fmt.Println(start.InitLogger())
	fmt.Println(start.InitAlert())
	//fmt.Println(start.InitAi())
	analyzer, err := ai.NewAiAnalyzer(&config2.Config.Ai)
	if err != nil {
		fmt.Println(err)
	}
	ai2.Registry(analyzer)
	tsk := ai2.NewTask(&config.AiTaskConfig{
		AiTaskName: "1",
		Cron: &config.Cron{
			Duration: time.Second * 10,
		},
		LogID: 1,
		LogFilter: config.LogFilter{
			StartTime: time.Now().AddDate(0, 0, -3),
			InspNames: []config.Name{"1"},
		},
		AlertID: 3,
	})
	fmt.Println(tsk.Do())
}

func TestAi(t *testing.T) {
	start.SetConfigPath("../../app/config", "yaml")
	start.Init()
	start.Run()
}
