package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/genai"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppBot struct {
	ctx       context.Context
	client    *whatsmeow.Client
	dbLog     waLog.Logger
	clientLog waLog.Logger
	genaiSvc  genai.CustomGenAIInterface
}

func NewWhatsAppBot(ctx context.Context) (*WhatsAppBot, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)

	address := fmt.Sprintf("file:%s?_foreign_keys=on", env.AppEnv.BotDBPath)
	storeContainer, err := sqlstore.New(ctx, "sqlite3", address, dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp store: %w", err)
	}

	deviceStore, err := storeContainer.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get WhatsApp device store: %w", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	bot := &WhatsAppBot{
		ctx:       ctx,
		client:    client,
		dbLog:     dbLog,
		clientLog: clientLog,
		genaiSvc:  genai.GenAI,
	}

	return bot, nil
}

func (s *WhatsAppBot) Start(ctx context.Context) error {
	s.clientLog.Infof("Starting WhatsApp bot...")

	s.client.AddEventHandler(s.eventHandler)

	if s.client.Store.ID == nil {
		qrChan, _ := s.client.GetQRChannel(ctx)
		err := s.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			}
		}

		return nil
	}

	err := s.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	return nil
}

func (s *WhatsAppBot) Stop() {
	if s.client != nil {
		s.clientLog.Infof("Disconnecting WhatsApp bot...")
		s.client.Disconnect()
	}
}

func (s *WhatsAppBot) eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		go s.handleMessage(v)
	case *events.Connected:
		s.clientLog.Infof("WhatsApp bot connected successfully")
	case *events.Disconnected:
		s.clientLog.Warnf("WhatsApp bot disconnected")
	case *events.LoggedOut:
		s.clientLog.Warnf("WhatsApp bot logged out. Please scan QR code again on next restart")
	case *events.StreamReplaced:
		s.clientLog.Warnf("WhatsApp bot stream replaced (logged in from another location)")
	}
}
