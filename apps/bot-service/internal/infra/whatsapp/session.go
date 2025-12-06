package whatsapp

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/types"
)

func getJakartaTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}

func getTimeBasedGreeting(t time.Time) string {
	hour := t.Hour()

	switch {
	case hour >= 4 && hour < 11:
		return "Selamat pagi"
	case hour >= 11 && hour < 15:
		return "Selamat siang"
	case hour >= 15 && hour < 18:
		return "Selamat sore"
	default:
		return "Selamat malam"
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

type sessionAction struct {
	actionType  string // "send_prompt", "auto_submit", "auto_close"
	phoneNumber string
	chatJID     types.JID
	message     string // only for send_prompt
}

func (s *WhatsAppBot) processExpiredSessions() {
	// First pass: collect actions while holding lock
	var actions []sessionAction

	s.sessionsMux.Lock()
	now := time.Now()
	for phoneNumber, session := range s.sessions {
		if session.WaitingForRating || session.WaitingForComment {
			continue
		}

		// Check if we need to send feedback prompt
		if !session.FeedbackPromptSent && now.Sub(session.LastMessageAt) > feedbackPromptDelay {
			// Update session state while locked
			session.FeedbackPromptSent = true
			session.IsAutoPrompt = true
			promptTime := now
			session.FeedbackPromptSentAt = &promptTime

			// Collect action to perform
			greeting := getTimeBasedGreeting(getJakartaTime())
			feedbackMessage := greeting + " Bapak/Ibu, untuk meningkatkan kualitas pelayanan kami, mohon dibantu penilaiannya 🙏🏻\n\nApabila berkenan, mohon kesediaan Bapak/Ibu untuk memberikan feedback terhadap kualitas pelayanan kami dengan rating 1-5.\n\nAdapun 3 poin penilaian sebagai berikut:\n1. Kecepatan dalam merespon pertanyaan/keluhan\n2. Kualitas komunikasi dan informasi yang diberikan\n3. Ketepatan dan kegunaan solusi yang diberikan\n\nUntuk memberikan feedback, silakan ketik /selesai\n\n⏱️ *Catatan:* Jika tidak ada respons dalam 5 menit, kami akan mencatat feedback Anda sebagai rating 5 bintang sebagai bentuk kepuasan terhadap layanan kami."

			actions = append(actions, sessionAction{
				actionType:  "send_prompt",
				phoneNumber: phoneNumber,
				chatJID:     *session.ChatJID,
				message:     feedbackMessage,
			})
		}

		// Check if session expired and needs auto-submit or close
		if session.FeedbackPromptSent && session.FeedbackPromptSentAt != nil {
			if now.Sub(*session.FeedbackPromptSentAt) > sessionExpiryDuration {
				if session.IsAutoPrompt {
					actions = append(actions, sessionAction{
						actionType:  "auto_submit",
						phoneNumber: phoneNumber,
						chatJID:     *session.ChatJID,
					})
				} else {
					actions = append(actions, sessionAction{
						actionType:  "auto_close",
						phoneNumber: phoneNumber,
					})
				}
			}
		}
	}
	s.sessionsMux.Unlock()

	// Execute actions without holding lock
	var sessionsToDelete []string
	for _, action := range actions {
		switch action.actionType {
		case "send_prompt":
			s.clientLog.Infof("Sending feedback prompt to %s due to inactivity", action.phoneNumber)
			s.sendMessage(action.chatJID, action.message)

		case "auto_submit":
			s.clientLog.Infof("Auto-submitting feedback rating 5 for %s due to no response", action.phoneNumber)
			s.autoSubmitFeedback(action.phoneNumber, action.chatJID)
			sessionsToDelete = append(sessionsToDelete, action.phoneNumber)

		case "auto_close":
			s.clientLog.Infof("Auto-closing session for %s due to no feedback response", action.phoneNumber)
			sessionsToDelete = append(sessionsToDelete, action.phoneNumber)
		}
	}

	// Delete sessions while holding lock
	if len(sessionsToDelete) > 0 {
		s.sessionsMux.Lock()
		for _, phoneNumber := range sessionsToDelete {
			delete(s.sessions, phoneNumber)
		}
		s.sessionsMux.Unlock()
	}
}

func (s *WhatsAppBot) getSession(phoneNumber string) *Session {
	s.sessionsMux.RLock()
	defer s.sessionsMux.RUnlock()

	return s.sessions[phoneNumber]
}

func (s *WhatsAppBot) createSession(phoneNumber string, chatJID *types.JID) *Session {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()

	session := &Session{
		PhoneNumber:   phoneNumber,
		LastMessageAt: time.Now(),
		ChatJID:       chatJID,
	}
	s.sessions[phoneNumber] = session

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
