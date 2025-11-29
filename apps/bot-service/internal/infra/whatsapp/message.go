package whatsapp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (s *WhatsAppBot) handleMessage(msg *events.Message) {
	if msg.Info.IsFromMe {
		return
	}

	meta := map[string]any{
		"pushname":  msg.Info.PushName,
		"timestamp": msg.Info.Timestamp,
	}

	if msg.Info.Type != "" {
		meta["type"] = msg.Info.Type
	}
	if msg.Info.Category != "" {
		meta["category"] = msg.Info.Category
	}
	if msg.IsViewOnce {
		meta["view_once"] = true
	}

	phoneNumber := msg.Info.Sender.User

	text := msg.Message.GetConversation()
	quotedMsg := ""
	if text == "" {
		text = msg.Message.GetExtendedTextMessage().GetText()
		quotedMsg = msg.Message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage().GetConversation()
	}

	if text == "" {
		return
	}

	log.Debug(log.CustomLogInfo{
		"from":      phoneNumber,
		"text":      text,
		"meta":      meta,
		"quotedMsg": quotedMsg,
	}, "[WhatsAppBot] Received WhatsApp message")

	s.statesMutex.RLock()
	userState, exists := s.userStates[phoneNumber]
	s.statesMutex.RUnlock()

	if exists && userState.WaitingForRating {
		s.handleRatingInput(msg, phoneNumber, text, userState)
		return
	}

	if exists && userState.WaitingForComment {
		s.handleCommentInput(msg, phoneNumber, text, userState)
		return
	}

	if strings.ToLower(strings.TrimSpace(text)) == "/selesai" {
		s.handleEndSession(msg, phoneNumber)
		return
	}

	session, err := s.sessionService.GetOrCreateSession(s.ctx, phoneNumber)
	if err != nil {
		s.clientLog.Errorf("Failed to get or create session: %v", err)
		s.sendReply(msg, "Sorry, something went wrong. Please try again later.")
		return
	}

	if err := s.sessionService.UpdateSessionActivity(s.ctx, session.ID); err != nil {
		s.clientLog.Errorf("Failed to update session activity: %v", err)
	}

	if strings.Contains(strings.ToLower(text), "sayang") {
		res, err := s.genaiSvc.Chat(s.ctx, []string{text})
		if err != nil {
			s.sendReply(msg, "Sorry, I couldn't process your message right now.")
			return
		}

		s.sendReply(msg, res)
	}
}

func (s *WhatsAppBot) handleEndSession(msg *events.Message, phoneNumber string) {
	session, err := s.sessionService.GetOrCreateSession(s.ctx, phoneNumber)
	if err != nil {
		s.clientLog.Errorf("Failed to get session: %v", err)
		s.sendReply(msg, "Sorry, something went wrong. Please try again later.")
		return
	}

	s.statesMutex.Lock()
	s.userStates[phoneNumber] = &UserState{
		SessionID:        session.ID.String(),
		WaitingForRating: true,
	}
	s.statesMutex.Unlock()

	s.sendReply(msg, "Terima kasih telah menggunakan layanan kami! 🙏\n\nSilakan berikan rating Anda (1-5):")
}

func (s *WhatsAppBot) handleRatingInput(msg *events.Message, phoneNumber string, text string, userState *UserState) {
	rating, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil || rating < 1 || rating > 5 {
		s.sendReply(msg, "Rating tidak valid. Silakan masukkan angka antara 1-5:")
		return
	}

	userState.Rating = rating
	userState.WaitingForRating = false
	userState.WaitingForComment = true

	s.statesMutex.Lock()
	s.userStates[phoneNumber] = userState
	s.statesMutex.Unlock()

	s.sendReply(msg, fmt.Sprintf("Terima kasih! Anda memberikan rating %d ⭐\n\nSilakan berikan komentar atau saran Anda (atau ketik '/skip' untuk melewati):", rating))
}

func (s *WhatsAppBot) handleCommentInput(msg *events.Message, phoneNumber string, text string, userState *UserState) {
	session, err := s.sessionService.GetOrCreateSession(s.ctx, phoneNumber)
	if err != nil {
		s.clientLog.Errorf("Failed to get session: %v", err)
		s.sendReply(msg, "Sorry, something went wrong. Please try again later.")
		return
	}

	var comment *string
	trimmedText := strings.TrimSpace(text)
	if strings.ToLower(trimmedText) != "/skip" && trimmedText != "" {
		comment = &trimmedText
	}

	feedbackID, err := uuid.UUID.NewV7()
	if err != nil {
		s.clientLog.Errorf("Failed to generate feedback ID: %v", err)
		s.sendReply(msg, "Sorry, something went wrong. Please try again later.")
		return
	}

	feedback := &entity.Feedback{
		ID:          feedbackID,
		SessionID:   session.ID,
		PhoneNumber: phoneNumber,
		Rating:      userState.Rating,
		Comment:     comment,
		CreatedAt:   time.Now(),
	}

	if err := s.feedbackRepo.Create(s.ctx, feedback); err != nil {
		s.clientLog.Errorf("Failed to save feedback: %v", err)
		s.sendReply(msg, "Sorry, something went wrong while saving your feedback.")
		return
	}

	s.statesMutex.Lock()
	delete(s.userStates, phoneNumber)
	s.statesMutex.Unlock()

	if err := s.sessionService.CloseSession(s.ctx, session.ID); err != nil {
		s.clientLog.Errorf("Failed to close session: %v", err)
	}

	s.sendReply(msg, "Terima kasih atas feedback Anda! 🙏\n\nSampai jumpa lagi! 👋")

	log.Info(log.CustomLogInfo{
		"phone_number": phoneNumber,
		"rating":       userState.Rating,
		"has_comment":  comment != nil,
		"feedback_id":  feedback.ID.String(),
	}, "[WhatsAppBot] Feedback received and saved")
}

func (s *WhatsAppBot) sendReply(msg *events.Message, text string) {
	_, err := s.client.SendMessage(s.ctx, msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(msg.Info.ID),
				Participant:   proto.String(msg.Info.Sender.String()),
				QuotedMessage: msg.Message,
			},
		},
	})
	if err != nil {
		s.clientLog.Errorf("Failed to send WhatsApp reply message: " + err.Error())
	}
}

func (s *WhatsAppBot) sendMessage(to types.JID, text string) {
	_, err := s.client.SendMessage(s.ctx, to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
		},
	})
	if err != nil {
		s.clientLog.Errorf("Failed to send WhatsApp message: " + err.Error())
	}
}
