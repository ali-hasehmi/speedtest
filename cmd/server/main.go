package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ali-hasehmi/speedtest/internal/config"
	"github.com/ali-hasehmi/speedtest/internal/handlers"
	"github.com/ali-hasehmi/speedtest/internal/metadata"
	"github.com/ali-hasehmi/speedtest/internal/speedtest"
	"github.com/ali-hasehmi/speedtest/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	majorVerion  = 0
	minorVersion = 1
	patchVersion = 0
	configPath   = "speedtest.env"
)

func VersionString() string {
	return fmt.Sprintf("%v.%v.%v", majorVerion, minorVersion, patchVersion)
}

func main() {

	logger.Infof("speedtest %v started (%v %v/%v)\n", VersionString(), runtime.Version(),
		runtime.GOOS, runtime.GOARCH)

	if err := config.Load(configPath); err != nil {
		logger.Fatal("failed to load config:", err)
	}

	logger.Infof(`Config loaded from %v with values: ListenAddr: %v, ListenPort: %v, ReadTimeout: %v, WriteTimeout: %v, IdleTimeout: %v, DownloadBufferSize: %v, DownloadMaxSize: %v, UploadMaxSize: %v, CityDBPath: %s, AsnDBPath: %s, LogFile: %s, LogLevel: %v`,
		configPath,
		config.ListenAddr(),
		config.ListenPort(),
		config.ReadTimeout(),
		config.WriteTimeout(),
		config.IdleTimeout(),
		config.DownloadBufferSize(),
		config.DownloadMaxSize(),
		config.UploadMaxSize(),
		config.CityDBPath(),
		config.AsnDBPath(),
		config.LogFile(),
		config.LogLevel(),
	)

	speedtest.InitBuffer(config.DownloadBufferSize())

	logger.SetLevel(config.LogLevel())

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	if err := metadata.Init(config.CityDBPath(), config.AsnDBPath()); err != nil {
		logger.Warningf("IP metadata partially disabled: %v", err)
	}
	defer metadata.Close()

	r.Get("/api/download", handlers.DownloadHandler)
	r.Post("/api/upload", handlers.UploadHandler)
	r.Get("/api/ip", handlers.IPHandler)
	r.Get("/api/ping", handlers.PingHandler)

	addr := fmt.Sprintf("%s:%d", config.ListenAddr(), config.ListenPort())

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
		IdleTimeout:  config.IdleTimeout(),
	}

	logger.Infof("server starting on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("listen: %s", err)
	}
}
