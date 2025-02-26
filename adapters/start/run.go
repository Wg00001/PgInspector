package start

import (
	"PgInspector/adapters/cron"
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

//todo:优雅关闭
