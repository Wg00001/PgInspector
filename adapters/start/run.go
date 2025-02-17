package start

import (
	"PgInspector/adapters/task/cron"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/17
 */

func Run() {
	cron.Start()
	select {}
}
