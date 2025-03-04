package start

import (
	"PgInspector/adapters/cron"
	"time"
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

func RunWithTimeAfter(duration time.Duration) {
	cron.Start()
	select {
	case <-time.After(duration):
		Close()
		return
	}
}

func Close() {
	cron.Start()
}

//todo:优雅关闭
