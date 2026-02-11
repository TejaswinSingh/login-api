package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/TejaswinSingh/login-api/internal/config"
	db "github.com/TejaswinSingh/login-api/internal/db/postgres"
	"github.com/TejaswinSingh/login-api/internal/logging"
	"github.com/TejaswinSingh/login-api/internal/metrics"
	"github.com/TejaswinSingh/login-api/internal/middleware"
	user "github.com/TejaswinSingh/login-api/internal/repository/user"
	"github.com/TejaswinSingh/login-api/internal/service/auth"
	"github.com/TejaswinSingh/login-api/internal/service/auth/jwt"
	userService "github.com/TejaswinSingh/login-api/internal/service/user"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := config.NewEnvConfig()

	logger := logging.NewLogger(config)

	logger.Info("ENV set to " + string(config.Env))

	dbpool, err := db.NewDbConnPool(config.DbConfig)
	if err != nil {
		logger.Error("unable to create db connection pool", "error", err.Error())
		panic("unable to create db connection pool")
	}
	defer dbpool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := dbpool.Ping(ctx); err != nil {
		logger.Error("unable to ping db", "error", err.Error())
		panic("unable to ping db")
	}

	userRepository := user.NewPostgresUserRepository(dbpool, logger, config)
	jwtModule := jwt.NewJwtModule(logger, config)
	metricsReg := metrics.NewMetricsRegistry()
	stdMetrics := metrics.NewStdMetrics(metricsReg)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(metricsReg, promhttp.HandlerOpts{}))

	// register service handlers
	auth.NewAuthService(logger, config, jwtModule, userRepository).RegisterHandlers(mux)
	userService.NewUserService(logger, config, userRepository).RegisterHandlers(mux)

	// wrap middlewares on mux
	handler := middleware.Metrics(mux, stdMetrics)
	handler = middleware.RequestLogger(handler, logger)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: handler,
	}

	shutdownDone := make(chan struct{})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan

		logger.Info("starting graceful shutdown for http server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("unable to complete graceful shutdown for http server", "error", err.Error())
		}
		close(shutdownDone)
	}()

	logger.Info("starting http server on port " + strconv.Itoa(config.HttpPort))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("http server stopped listening", "error", err.Error())
		os.Exit(1)
	}

	<-shutdownDone
	logger.Info("exiting program")
}
