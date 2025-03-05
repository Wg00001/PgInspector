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

func Run() func() {
	ch := make(chan struct{})

	go func() {
		cron.Start()
		select {
		case <-ch:
			cron.Exit()
			break
		}
	}()

	return func() {
		close(ch)
	}
}

func RunWithTimeAfter(duration time.Duration) {
	cron.Start()
	select {
	case <-time.After(duration):
		cron.Exit()
		return
	}
}
