package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	feedbackrepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/repository"
	userrepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/repository"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/genai"
	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppBot struct {
	ctx          context.Context
	client       *whatsmeow.Client
	dbLog        waLog.Logger
	clientLog    waLog.Logger
	genaiSvc     genai.CustomGenAIInterface
	feedbackRepo contracts.FeedbackRepository
	userRepo     contracts.UserRepository
	sessions     map[string]*Session
	sessionsMux  sync.RWMutex
}

type Session struct {
	PhoneNumber          string
	LastMessageAt        time.Time
	WaitingForRating     bool
	WaitingForComment    bool
	Rating               int
	FeedbackPromptSent   bool
	FeedbackPromptSentAt *time.Time
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

	feedbackRepo := feedbackrepository.NewFeedbackRepository(sqlxDB)
	userRepo := userrepository.NewUserRepository(sqlxDB)

	bot := &WhatsAppBot{
		ctx:          ctx,
		client:       client,
		dbLog:        dbLog,
		clientLog:    clientLog,
		genaiSvc:     genai.GenAI,
		feedbackRepo: feedbackRepo,
		userRepo:     userRepo,
		sessions:     make(map[string]*Session),
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
			s.processExpiredSessions()
		}
	}
}

func (s *WhatsAppBot) processExpiredSessions() {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()

	now := time.Now()
	for phoneNumber, session := range s.sessions {
		if session.WaitingForRating || session.WaitingForComment {
			continue
		}

		if !session.FeedbackPromptSent && now.Sub(session.LastMessageAt) > 5*time.Minute {
			session.FeedbackPromptSent = true
			promptTime := now
			session.FeedbackPromptSentAt = &promptTime
			s.clientLog.Infof("Sending feedback prompt to %s due to inactivity", phoneNumber)
		}

		if session.FeedbackPromptSent && session.FeedbackPromptSentAt != nil {
			if now.Sub(*session.FeedbackPromptSentAt) > 5*time.Minute {
				s.clientLog.Infof("Auto-closing session for %s due to no feedback response", phoneNumber)
				delete(s.sessions, phoneNumber)
			}
		}
	}
}

func (s *WhatsAppBot) getOrCreateSession(phoneNumber string) *Session {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()

	session, exists := s.sessions[phoneNumber]
	if !exists {
		session = &Session{
			PhoneNumber:   phoneNumber,
			LastMessageAt: time.Now(),
		}
		s.sessions[phoneNumber] = session
	}

	return session
}

func (s *WhatsAppBot) updateSessionActivity(phoneNumber string) {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()

	if session, exists := s.sessions[phoneNumber]; exists {
		session.LastMessageAt = time.Now()
	}
}

func (s *WhatsAppBot) deleteSession(phoneNumber string) {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()

	delete(s.sessions, phoneNumber)
}
