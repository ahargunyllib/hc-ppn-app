package main

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/database"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/server"
)

func main() {
	server := server.NewHTTPServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	server.MountMiddlewares()
	server.MountRoutes(psqlDB)
	server.Start(env.AppEnv.AppPort)
}
