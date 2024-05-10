package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"os/signal"
	"service_template/internal/api"
	"service_template/internal/config"
	"service_template/internal/repository"
	"service_template/internal/repository/model"
	"service_template/internal/server"
	"service_template/internal/service"
	"service_template/pkg/cache"
	"service_template/pkg/db"
	"service_template/pkg/logger"
	"syscall"
	"time"
)

var (
	configFile = flag.String("config", "./config.yaml", "config file path")
	debugMode  = flag.Bool("debug", false, "open debug mode")
	migrate    = flag.Bool("migrate", false, "migrate tables")
	testData   = flag.Bool("test-data", false, "create test data to database")
)

func main() {
	flag.Parse()

	// 初始化配置
	cfg, err := config.InitConfig(*configFile)
	if err != nil {
		panic(err)
	}
	// 初始化日志
	lj := logger.NewLumberjack(cfg.Log.Lumberjack)
	cfg.Log.Output = []io.Writer{lj}
	cfg.Log.ErrOutput = []io.Writer{lj}
	if *debugMode {
		cfg.Log.Output = append(cfg.Log.Output, os.Stdout)
		cfg.Log.ErrOutput = append(cfg.Log.Output, os.Stderr)
	}
	log := logger.InitLogger(cfg.Log)
	// 初始化数据库
	cfg.Database.Debug = *debugMode
	database, err := db.NewDB(cfg.Database)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	if *migrate {
		err := database.AutoMigrate(model.Tables...)
		if err != nil {
			panic(err)
		}
		fmt.Println("Migrate tables finished.")
		return
	}
	if *testData {
		err := genTestData(database)
		if err != nil {
			panic(err)
		}
		fmt.Println("Generating test data finished.")
		return
	}
	// 初始化存储层
	repo := repository.NewRepository(database)
	// 初始化 Redis
	rdb, err := cache.NewRedis(cfg.Cache)
	if err != nil {
		return
	}
	defer rdb.Close()
	// 初始化服务层
	srv := service.NewService(repo, rdb)
	// 初始化接口层
	httpApi := api.InitApi(srv)
	if *debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// 启动服务
	engine := gin.New()
	engine.Use(gin.RecoveryWithWriter(log.GetOutput()))
	pprof.Register(engine, "/inner/"+pprof.DefaultPrefix)
	s := server.NewServer(engine, cfg.HttpServer)
	api.InitRouter(engine, httpApi, srv)
	api.RegisterSwagger(engine, s.Addr())
	done := make(chan struct{})
	go func() {
		err := s.Start()
		if err != nil {
			logger.Errorf("http server exited, error: %v", err)
		}
		done <- struct{}{}
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signals:
	case <-done:
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}

func genTestData(db *db.DB) error {
	return nil
}
