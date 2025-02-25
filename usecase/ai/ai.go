package ai

import (
	"PgInspector/entities/config"
	"PgInspector/usecase/logger"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

func Init(cfg config.AiConfig) {}

type AiTask config.AiTaskConfig

func NewTask(taskConfig config.AiTaskConfig) AiTask {
	return AiTask(taskConfig)
}

func (t *AiTask) Do() {
	contents, err := logger.Get(t.LogID).ReadLog(t.LogFilter)
	if err != nil {
		return
	}
}
