package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/handlers/sadtalker"
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
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
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

	handlers.BaseHdl.Setup(&handlers.BaseHandlerCfg{
		DB: db,
		Q:  taskQueue,
		L:  logger,
		C:  initClientSet(),
	})
	handlers.MidHdl.Setup(&handlers.MiddlewareHandlerCfg{L: logger, TicketExpire: false})
	sadtalker.StHdl.Setup(&sadtalker.Cfg{JobNamespace: cfg.SadTalker.JobNamespace})

	db.AutoMigrate(&models.Task{})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Metrics: metricsserver.Options{BindAddress: "0"},
	})
	if err != nil {
		logger.Fatal(err)
	}

	go StartWorker(taskQueue)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err = mgr.Start(ctx); err != nil {
			logger.Fatalf("%v", err)
		}
	}()

	router := gin.Default()
	routes.RegisterRoutes(router)
	if err = router.Run(":8080"); err != nil {
		panic(err)
	}
	cancel()
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
