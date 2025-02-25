package ai

import (
	"PgInspector/entities/config"
	"PgInspector/entities/task"
	"PgInspector/usecase/ai/format"
	"PgInspector/usecase/logger"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

type AiTask config.AiTaskConfig

func NewTask(taskConfig config.AiTaskConfig) AiTask {
	return AiTask(taskConfig)
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
		return fmt.Errorf("Ai task err\n- AiTask name: %v\n- err: log read empty\n---\n", t.AiTaskName)
	}
	//第三和第四步通过调用Analyzer完成
	//3. 发送Ai
	//4. 解析结果

	//5. 发送Alert
	return nil
}

func (t *AiTask) GetCron() *config.Cron {
	return t.Cron
}

func (t *AiTask) GetName() config.Name {
	return t.AiTaskName
}
