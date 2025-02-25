package cron

import (
	"PgInspector/entities/config"
	"PgInspector/entities/task"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"log"
	"strconv"
	"sync"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/14
 */

var (
	s  gocron.Scheduler
	mu sync.Mutex
)

func Init() {
	mu.Lock()
	defer mu.Unlock()
	sTemp, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local), // 设置时区
		gocron.WithGlobalJobOptions(),   // 全局任务选项
	)
	if err != nil {
		log.Printf("initScheduler失败！: %v", err)
		panic("initScheduler失败！: " + err.Error())
	}
	s = sTemp
	return
}

func AddTask(task *task.Task) {
	mu.Lock()
	defer mu.Unlock()
	definition, err := jobDefinition(task.Config.Cron)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = s.NewJob(
		definition,
		gocron.NewTask(func() {
			err := task.Do()
			if err != nil {
				log.Printf("gocron do task err: %v", err)
				return
			}
		}), // 任务函数和参数
		gocron.WithName(task.Config.TaskName.Str()))
	if err != nil {
		log.Println(err)
	}
}

func AddAiTask() {
	mu.Lock()
	defer mu.Unlock()

}

func Start() {
	s.Start()
	log.Println("cron start...")
}

func Monitor() {

}

// 将task中的时间设置读取到cron的对象中
func jobDefinition(t *config.Cron) (gocron.JobDefinition, error) {
	if t == nil {
		return nil, fmt.Errorf("gocron add task err, time not define, taskname: %s\n", t)
	}
	cConfig := *t

	if cConfig.Duration != 0 {
		return gocron.DurationJob(cConfig.Duration), nil
	}
	atTime := gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))
	if len(cConfig.AtTime) != 0 {
		atTime = gocron.NewAtTimes(
			gocron.NewAtTime(parseAtTime(cConfig.AtTime[0])),
			func() []gocron.AtTime {
				res := make([]gocron.AtTime, 0, len(cConfig.AtTime)-1)
				for i := 1; i < len(cConfig.AtTime); i++ {
					res = append(res, gocron.NewAtTime(parseAtTime(cConfig.AtTime[i])))
				}
				return res
			}()...)
	}
	if len(cConfig.Monthly) != 0 {
		return gocron.MonthlyJob(1, func() gocron.DaysOfTheMonth {
			if len(cConfig.Monthly) > 1 {
				return gocron.NewDaysOfTheMonth(cConfig.Monthly[0], cConfig.Monthly[1:]...)
			}
			return gocron.NewDaysOfTheMonth(cConfig.Monthly[0])
		}(), atTime), nil
	}
	if len(cConfig.Weekly) != 0 {
		return gocron.WeeklyJob(1, func() gocron.Weekdays {
			if len(cConfig.Weekly) > 1 {
				return gocron.NewWeekdays(cConfig.Weekly[0], cConfig.Weekly[1:]...)
			}
			return gocron.NewWeekdays(cConfig.Weekly[0])
		}(), atTime), nil
	}
	return gocron.DailyJob(1, atTime), nil
	//return gocron.DurationJob(time.Second * 5)
}

func parseAtTime(t string) (uint, uint, uint) {
	temp := [3]uint64{}
	for l, i := 0, 0; l < len(t) && i < 3; {
		if t[l] < '0' || t[l] > '9' {
			l++
			continue
		}
		r := l + 1
		for r < len(t) && t[r] >= '0' && t[r] <= '9' {
			r++
		}
		temp[i], _ = strconv.ParseUint(t[l:r], 10, 64)
		i++
		l = r
	}
	return uint(temp[0]), uint(temp[1]), uint(temp[2])
}
