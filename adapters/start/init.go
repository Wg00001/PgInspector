package start

import (
	"PgInspector/adapters/cron"
	"PgInspector/usecase/ai"
	"PgInspector/usecase/alerter"
	config2 "PgInspector/usecase/config"
	"PgInspector/usecase/db"
	"PgInspector/usecase/logger"
	"PgInspector/usecase/task"
	"fmt"
	"github.com/wg00001/wgo-sdk/wg"
	"log"
)
import (
	_ "PgInspector/adapters/ai/agent"
	_ "PgInspector/adapters/ai/default"
	_ "PgInspector/adapters/ai/ollama"
	_ "PgInspector/adapters/alerter/default"
	_ "PgInspector/adapters/alerter/empty"
	_ "PgInspector/adapters/alerter/feishu"
	_ "PgInspector/adapters/config/yaml"
	_ "PgInspector/adapters/logger/default"
	_ "PgInspector/adapters/logger/postgres"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/17
 */

func Init() {
	log.SetFlags(log.LstdFlags)

	//config.InitConfig(config.NewReader(pConfigType, pFilePath))
	err := config2.Open(pConfigType, map[string]string{
		"filepath": pFilePath,
	})
	if err != nil {
		panic(fmt.Sprintf("config open fail: %s", err))
	}
	err = InitDB()
	if err != nil {
		panic(fmt.Sprintf("db init fail: %s", err))
	}

	defer func() {
		if r := recover(); r != nil {
			db.CloseAll()
			panic(r)
		}
	}()
	printErr := func(err error) {
		if err != nil {
			panic(fmt.Sprintf("!!!!! System init fail !!!!!\n!!!!! Err :%s\n\n", err))
		}
	}

	cron.Init()
	printErr(InitLogger())
	printErr(InitTask())
	printErr(InitAlert())
	printErr(InitAiConfig())
	printErr(InitAiTask())
	log.Println("====== System NewReader Completely ======")
}

func InitDB() error {
	config2.RLock()
	defer config2.RUnlock()
	dbConfigs := wg.MapToValueSlice(config2.Config.DB)
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

func InitLogger() error {
	config2.RLock()
	defer config2.RUnlock()
	logConfigs := wg.MapToValueSlice(config2.Config.Log)
	for _, v := range logConfigs {
		err := logger.Use(*v)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitTask() error {
	config2.RLock()
	defer config2.RUnlock()
	taskConfigs := wg.MapToValueSlice(config2.Config.Task)
	for _, v := range taskConfigs {
		t, err := task.NewTask(v)
		if err != nil {
			return err
		}
		err = task.Register(t)
		if err != nil {
			return err
		}
		cron.AddTask(t)
	}
	return nil
}

func initCron() error {
	cron.Init()
	//taskConfigs := wg.MapToValueSlice(usecase.Config.Task)
	//for _, cfg := range taskConfigs {
	//	cron.AddTask(task.Get(cfg.TaskName))
	//}
	return nil
}

func InitAlert() error {
	config2.RLock()
	defer config2.RUnlock()
	alertConfigs := wg.MapToValueSlice(config2.Config.Alert)
	for _, v := range alertConfigs {
		//a, err := alerter.NewAlerter(v)
		//if err != nil {
		//	return err
		//}
		//err = alerter.Register(v.AlertID, a)
		//if err != nil {
		//	return err
		//}

		err := alerter.Use(*v)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitAiConfig() error {
	config2.RLock()
	defer config2.RUnlock()
	return ai.Use(config2.Config.Ai)
}

func InitAiTask() error {
	config2.RLock()
	defer config2.RUnlock()

	aiTasks := wg.MapToValueSlice(config2.Config.AiTask)
	for _, v := range aiTasks {
		cron.AddTask(ai.NewTask(v))
	}
	return nil
}
