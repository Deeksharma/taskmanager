package main

import (
	"context"
	"flag"
	"github.com/Deeksharma/taskmanager/internal/config"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/server"
	"github.com/Deeksharma/taskmanager/internal/server/api"
)

const (
	ApiServer = "api"
)

func main() {
	serverType := *flag.String("type", ApiServer, "server type")
	allConfigs := config.AllSettings()
	log.InfoWithFields(context.TODO(), map[string]interface{}{
		"config": allConfigs,
	}, "all configs")

	var server server.Server
	switch serverType {
	case ApiServer:
		server = api.NewApiServer()
	default:
		server = api.NewApiServer()
	}
	server.Start()

	wait := server.Stop()
	<-wait
}
