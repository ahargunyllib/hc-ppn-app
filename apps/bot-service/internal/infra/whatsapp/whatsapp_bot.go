package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	feedbackRepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/repository"
	feedbackService "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/service"
	userRepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/repository"
	userService "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/genai"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var (
	sessionExpiryDuration = 1 * time.Minute // Duration after which inactive sessions are cleared
	feedbackPromptDelay   = 1 * time.Minute // Duration of inactivity after which feedback prompt is sent
)

type WhatsAppBot struct {
	ctx         context.Context
	client      *whatsmeow.Client
	dbLog       waLog.Logger
	clientLog   waLog.Logger
	genaiSvc    genai.CustomGenAIInterface
	feedbackSvc contracts.FeedbackService
	userSvc     contracts.UserService
	sessions    map[string]*Session
	sessionsMux sync.RWMutex
}

type Session struct {
	PhoneNumber          string
	LastMessageAt        time.Time
	WaitingForRating     bool
	WaitingForComment    bool
	Rating               int
	FeedbackPromptSent   bool
	FeedbackPromptSentAt *time.Time
	ChatJID                  *types.JID
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

	validator := validator.Validator
	uuid := uuid.UUID

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	feedbackRepo := feedbackRepository.NewFeedbackRepository(sqlxDB)
	userRepo := userRepository.NewUserRepository(sqlxDB)

	feedbackSvc := feedbackService.NewFeedbackService(feedbackRepo, validator, uuid)
	userSvc := userService.NewUserService(userRepo, validator, uuid)

	bot := &WhatsAppBot{
		ctx:         ctx,
		client:      client,
		dbLog:       dbLog,
		clientLog:   clientLog,
		genaiSvc:    genai.GenAI,
		feedbackSvc: feedbackSvc,
		userSvc:     userSvc,
		sessions:    make(map[string]*Session),
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
