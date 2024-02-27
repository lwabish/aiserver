package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/handlers/roop"
	"github.com/lwabish/cloudnative-ai-server/handlers/sadtalker"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/routes"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
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

	var db *gorm.DB
	switch cfg.Db.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open("task.db"), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Db.Mysql), &gorm.Config{})
	}
	if err != nil {
		logger.Fatal(err)
	}
	if err = db.AutoMigrate(&models.Task{}); err != nil {
		logger.Fatal(err)
	}

	taskQueue := utils.NewTaskQueue()

	handlers.BaseHdl.Setup(&handlers.BaseHandlerCfg{
		DB: db,
		Q:  taskQueue,
		L:  logger,
	})
	handlers.MidHdl.Setup(&handlers.MiddlewareHandlerCfg{L: logger, TicketExpire: cfg.Auth.TokenExpire})
	sadtalker.StHdl.Setup(&cfg)
	roop.Handler.Setup(&cfg)

	go StartWorker(taskQueue)

	ctx, cancel := context.WithCancel(context.Background())
	if cfg.Mode == "cloud-native" {
		go enableController(ctx, logger)
		handlers.BaseHdl.SetupCloudNative(&handlers.BaseHandlerCfg{
			C: initClientSet(),
		})
	}

	router := gin.Default()
	routes.RegisterRoutes(router)
	// fixme: no panic
	if err = router.Run(":8080"); err != nil {
		panic(err)
	}
	cancel()
}

func enableController(ctx context.Context, logger *logrus.Logger) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Metrics: metricsserver.Options{BindAddress: "0"},
		Logger:  klog.NewKlogr(),
	})
	if err != nil {
		logger.Fatal(err)
	}

	err = (&controllers.SadTalkerJobReconciler{
		Client:      mgr.GetClient(),
		BaseHandler: handlers.BaseHdl,
		Logger:      mgr.GetLogger().WithName("sad-talker-reconciler"),
	}).SetupWithManager(mgr)
	if err != nil {
		logger.Fatalln(err)
	}

	if err = mgr.Start(ctx); err != nil {
		logger.Fatalf("%v", err)
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
