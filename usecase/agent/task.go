package agent

import (
	alerter2 "PgInspector/entities/alerter"
	"PgInspector/entities/config"
	"PgInspector/entities/task"
	"PgInspector/usecase/agent/format"
	"PgInspector/usecase/alerter"
	"PgInspector/usecase/logger"
	"encoding/json"
	"fmt"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

type AiTask struct {
	*config.AgentTaskConfig
}

func NewTask(taskConfig *config.AgentTaskConfig) *AiTask {
	return &AiTask{AgentTaskConfig: taskConfig}
}

var _ task.Task = (*AiTask)(nil)

func (t *AiTask) Do() error {
	//1. 获取日志
	contents, err := logger.Get(t.LogID).ReadLog(t.LogFilter)
	if err != nil {
		return err
	}
	//2. 组织格式
	msg, err := format.Format(contents...)
	if err != nil {
		return err
	}
	if msg == nil {
		return fmt.Errorf("Ai task err\n- AiTask name: %v\n- err: log read empty\n---\n", t.Name)
	}
	//第三和第四步通过调用Analyzer完成

	//3. 发送Ai
	//4. 解析结果
	res, err := Analyze(*msg)
	if err != nil {
		return err
	}

	//5. 发送Alert
	return alerter.GetAlert(t.AlertID).Send(*buildAiAlertContent(t, res))
}

func (t *AiTask) GetCron() *config.Cron {
	return t.Cron
}

func (t *AiTask) GetName() config.Name {
	return "ai_task:" + t.Name
}

func buildAiAlertContent(t *AiTask, msg string) *alerter2.Content {
	return &alerter2.Content{
		TimeStamp: time.Now(),
		TaskName:  "AiAnalyzeTask: " + t.Name,
		TaskID:    time.Now().Format("20060504_150201"),
		InspName:  logFilterString(t.LogFilter),
		Result:    []map[string]interface{}{{"message": msg}},
	}
}

func logFilterString(filter config.LogFilter) string {
	marshal, err := json.Marshal(filter)
	if err != nil {
		return ""
	}
	return string(marshal)
}
