package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"github.com/lwabish/cloudnative-ai-server/controllers/sadtalker"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/routes"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path"
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

	taskQueue := utils.NewTaskQueue()

	controllers.BaseCtl.Setup(&controllers.BaseControllerCfg{
		DB: db,
		Q:  taskQueue,
		L:  logger,
		C:  initClientSet(),
	})
	controllers.MidCtl.Setup(&controllers.MiddlewareControllerCfg{L: logger, TicketExpire: false})
	sadtalker.StCtl.Setup(&sadtalker.Cfg{JobNamespace: cfg.SadTalker.JobNamespace})

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

func initClientSet() *kubernetes.Clientset {
	var c *rest.Config
	var err error
	c, err = clientcmd.BuildConfigFromFlags("", path.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		c, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}
	client, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}
