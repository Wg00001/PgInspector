package start

import (
	"PgInspector/adapters/config_adapter"
	"PgInspector/adapters/logger_adapter"
	"PgInspector/adapters/task/cron"
	"PgInspector/entities/config"
	"PgInspector/usecase"
	"PgInspector/usecase/db"
	"PgInspector/usecase/logger"
	"PgInspector/usecase/task"
	"fmt"
	"github.com/wg00001/wgo-sdk/wg"
	"log"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/17
 */

func Init() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	log.SetFlags(log.LstdFlags)

	config.InitConfig(config_adapter.NewReader(pConfigType, pFilePath))
	err := initDB()
	if err != nil {
		log.Printf("db init fail: %s", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			db.CloseAll()
		}
	}()
	printErr := func(err error) {
		if err != nil {
			panic(fmt.Sprintf("System init fail :\n Err :%s\n", err))
		}
	}

	printErr(initLogger())
	printErr(initTask())
	printErr(initCron())
	log.Println("System Init Succeed !!!")
}

func initDB() error {
	usecase.RLock()
	defer usecase.RUnlock()
	dbConfigs := wg.MapToValueSlice(usecase.Config.DB)
	for _, v := range dbConfigs {
		sqlDB, err := db.InitDB(v)
		if err != nil {
			return err
		}
		err = db.Register(sqlDB)
		if err != nil {
			return err
		}
	}
	return nil
}

func initLogger() error {
	usecase.RLock()
	defer usecase.RUnlock()
	logConfigs := wg.MapToValueSlice(usecase.Config.Log)
	for _, v := range logConfigs {
		l, err := logger_adapter.NewLogger(v)
		if err != nil {
			return err
		}
		err = logger.Register(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func initTask() error {
	usecase.RLock()
	defer usecase.RUnlock()
	taskConfigs := wg.MapToValueSlice(usecase.Config.Task)
	for _, v := range taskConfigs {
		t, err := task.InitTask(v)
		if err != nil {
			return err
		}
		err = task.Register(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func initCron() error {
	cron.Init()
	taskConfigs := wg.MapToValueSlice(usecase.Config.Task)
	for _, cfg := range taskConfigs {
		cron.AddTask(task.Get(cfg.TaskName))
	}
	return nil
}
