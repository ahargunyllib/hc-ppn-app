package whatsapp

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/types"
)

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

		if !session.FeedbackPromptSent && now.Sub(session.LastMessageAt) > feedbackPromptDelay {
			session.FeedbackPromptSent = true
			promptTime := now
			session.FeedbackPromptSentAt = &promptTime
			s.clientLog.Infof("Sending feedback prompt to %s due to inactivity", phoneNumber)
			s.sendMessage(*session.ChatJID, "Sudah cukup lama sejak kami menerima pesan dari Anda. Mohon berikan feedback tentang layanan kami. Untuk mengirim feedback, ketik /selesai.")
		}

		if session.FeedbackPromptSent && session.FeedbackPromptSentAt != nil {
			if now.Sub(*session.FeedbackPromptSentAt) > sessionExpiryDuration {
				s.clientLog.Infof("Auto-closing session for %s due to no feedback response", phoneNumber)
				delete(s.sessions, phoneNumber)
			}
		}
	}
}

func (s *WhatsAppBot) getOrCreateSession(phoneNumber string, chatJID *types.JID) *Session {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()

	session, exists := s.sessions[phoneNumber]
	if !exists {
		session = &Session{
			PhoneNumber:   phoneNumber,
			LastMessageAt: time.Now(),
			ChatJID:           chatJID,
		}
		s.sessions[phoneNumber] = session

		s.sendMessage(*chatJID, "Halo! Selamat datang di layanan WhatsApp kami. Jika anda sudah selesai, silakan ketik /selesai untuk memberikan feedback.")
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
