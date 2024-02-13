package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/routes"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		os.Exit(1)
	}
	logger := utils.NewLogger(level)

	db, err := gorm.Open("mysql", cfg.DatabaseURL)
	defer func(db *gorm.DB) {
		logger.Fatal(db.Close())
	}(db)
	if err != nil {
		logger.Fatal(err)
	}

	// 初始化任务队列
	taskQueue := utils.NewTaskQueue(0)

	controllers.Inject(&controllers.BaseControllerCfg{
		DB: db,
		Q:  taskQueue,
		L:  logger,
	})

	db.AutoMigrate(&models.Task{})

	// 启动工作goroutine
	go StartWorker(taskQueue)

	// 启动http server
	router := gin.Default()
	routes.RegisterRoutes(router)
	if err = router.Run(":8080"); err != nil {
		panic(err)
	}
}
