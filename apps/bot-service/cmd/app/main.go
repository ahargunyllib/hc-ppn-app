package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/database"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/server"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/whatsapp"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
)

func main() {
	server := server.NewHTTPServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	server.MountMiddlewares()
	server.MountRoutes(psqlDB)

	if env.AppEnv.BotEnabled {
		go startWhatsAppBot()
	}

	server.Start(env.AppEnv.AppPort)
}

func startWhatsAppBot() {
	ctx := context.Background()

	botService, err := whatsapp.NewWhatsAppBot(ctx)
	if err != nil {
		log.Error(log.CustomLogInfo{
			"error": err.Error(),
		}, "[WhatsAppBot] Failed to create WhatsApp service")
		return
	}

	if err := botService.Start(ctx); err != nil {
		log.Error(log.CustomLogInfo{
			"error": err.Error(),
		}, "[WhatsAppBot] Failed to start WhatsApp service")
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	botService.Stop()
}
