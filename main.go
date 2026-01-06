package main

import (
	"app/api"
	"app/lib"
	"app/lib/config"
	"app/lib/ws"
	"app/middleware"
	"app/repository/dao"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func setupApp() *gin.Engine {
	apiLogger := lib.NewLogger(filepath.Join(config.App.LogDir, "api.log"))
	appLogger := lib.NewLogger(filepath.Join(config.App.LogDir, "app.log"))
	defer apiLogger.Sync()
	defer appLogger.Sync()
	app := gin.New()
	app.Use(middleware.Logger(apiLogger))
	app.Use(middleware.Recovery(appLogger))
	app.Use(middleware.Error())
	app.Use(middleware.Cors())
	lib.InitTranslator(config.App.Locale)
	lib.RegisterValidatorTranslations(config.App.Locale)
	dao.Init(config.App.Dsn)
	api.ApplyRoutes(app)
	go ws.WebsocketManager.Start()
	return app
}

var version = ""

var printVersion bool

func init() {
	flag.BoolVar(&printVersion, "version", false, "print program build version")
	flag.Parse()
}

func main() {
	if printVersion {
		println(version)
		os.Exit(0)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config.Read()
	app := setupApp()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.App.Port),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	stop()
	log.Println("shutdown gracefully, press ctrl+c force shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("failed to shutdown server: ", err)
	}
	log.Println("server exiting")
}
