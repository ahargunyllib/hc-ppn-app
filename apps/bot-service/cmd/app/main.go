package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/database"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/server"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/whatsapp"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	server := server.NewHTTPServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	server.MountMiddlewares()
	server.MountRoutes(psqlDB)

	if env.AppEnv.BotEnabled {
		wg.Add(1)
		go startWhatsAppBot(ctx, psqlDB.DB, &wg)
	}

	go server.Start(env.AppEnv.AppPort)

	<-ctx.Done()
	log.Info(log.CustomLogInfo{}, "Shutting down gracefully...")

	// Wait for all goroutines to finish cleanup
	wg.Wait()
	log.Info(log.CustomLogInfo{}, "Shutdown complete")
}

func startWhatsAppBot(ctx context.Context, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	botService, err := whatsapp.NewWhatsAppBot(ctx, db)
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

	<-ctx.Done()
	botService.Stop()
	log.Info(log.CustomLogInfo{}, "WhatsApp service stopped")
}
