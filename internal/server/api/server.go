package api

import (
	"context"
	"errors"
	"github.com/Deeksharma/taskmanager/internal/config"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/repository/dbrepo"
	"github.com/Deeksharma/taskmanager/internal/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ApiServer struct { // should have gin object
	srv               *http.Server
	ctx               context.Context
	closeDBConnection func() error
}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

func (as *ApiServer) Start() {
	ctx := context.Background()
	as.ctx = ctx
	log.Info(ctx, "Starting Api Server...")
	taskDatabaseRepo, closeDConnection := dbrepo.NewDBRepo(ctx)
	as.closeDBConnection = closeDConnection

	router := routes.NewRouter(taskDatabaseRepo)
	srv := &http.Server{
		Addr:    config.GetString("server.host") + ":" + config.GetString("server.port"),
		Handler: router,
	}
	as.srv = srv
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(ctx, "Failed to initialize server: %v\n", err)
		}
	}()
	log.Infof(ctx, "Listening on port %v", srv.Addr)

}

func (as *ApiServer) Stop() <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
		<-s

		log.Info(as.ctx, "shutting down...")

		// set timeout for the ops to be done to prevent system hang
		timeout := 20 * time.Second
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Infof(as.ctx, "timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		log.Info(as.ctx, "shutting down database")
		if err := as.closeDBConnection(); err != nil {
			log.Fatalf(as.ctx, "Database forced to shutdown: %v\n", err)
		}
		log.Info(as.ctx, "shutting down server")
		if err := as.srv.Shutdown(as.ctx); err != nil {
			log.Fatalf(as.ctx, "Server forced to shutdown: %v\n", err)
		}
		close(wait)
	}()

	return wait
}
