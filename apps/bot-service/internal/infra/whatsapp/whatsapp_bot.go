package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	feedbackrepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/repository"
	sessionrepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/session/repository"
	sessionservice "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/session/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/genai"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppBot struct {
	ctx            context.Context
	client         *whatsmeow.Client
	dbLog          waLog.Logger
	clientLog      waLog.Logger
	genaiSvc       genai.CustomGenAIInterface
	sessionService *sessionservice.SessionService
	feedbackRepo   contracts.FeedbackRepository
	userStates     map[string]*UserState
	statesMutex    sync.RWMutex
}

type UserState struct {
	SessionID         string
	WaitingForRating  bool
	WaitingForComment bool
	Rating            int
}

func NewWhatsAppBot(ctx context.Context, db *sql.DB, sqlxDB *sqlx.DB) (*WhatsAppBot, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)

	storeContainer := sqlstore.NewWithDB(db, "postgres", dbLog)
	err := storeContainer.Upgrade(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade WhatsApp database store: %w", err)
	}

	deviceStore, err := storeContainer.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get WhatsApp device store: %w", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	sessionRepo := sessionrepository.NewSessionRepository(sqlxDB)
	feedbackRepo := feedbackrepository.NewFeedbackRepository(sqlxDB)
	logger := log.NewLogger()
	sessionService := sessionservice.NewSessionService(sessionRepo, uuid.UUID, logger)

	bot := &WhatsAppBot{
		ctx:            ctx,
		client:         client,
		dbLog:          dbLog,
		clientLog:      clientLog,
		genaiSvc:       genai.GenAI,
		sessionService: sessionService,
		feedbackRepo:   feedbackRepo,
		userStates:     make(map[string]*UserState),
	}

	return bot, nil
}

func (s *WhatsAppBot) Start(ctx context.Context) error {
	s.clientLog.Infof("Starting WhatsApp bot...")

	s.client.AddEventHandler(s.eventHandler)

	go s.sessionExpiryChecker(ctx)

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

func (s *WhatsAppBot) sessionExpiryChecker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.sessionService.ProcessExpiredSessions(ctx); err != nil {
				s.clientLog.Errorf("Failed to process expired sessions: %v", err)
			}
		}
	}
}
