package task

import "PgInspector/entities/config"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

// Task 任务接口，巡检任务和ai分析任务实现此接口，用于交给cron分析。
type Task interface {
	Do() error
	GetCron() *config.Cron
	GetName() config.Name
}
