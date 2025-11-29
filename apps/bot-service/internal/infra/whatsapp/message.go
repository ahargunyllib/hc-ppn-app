package whatsapp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
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

	session := s.getOrCreateSession(phoneNumber)

	if session.WaitingForRating {
		s.handleRatingInput(msg, phoneNumber, text, session)
		return
	}

	if session.WaitingForComment {
		s.handleCommentInput(msg, phoneNumber, text, session)
		return
	}

	if strings.ToLower(strings.TrimSpace(text)) == "/selesai" {
		s.handleEndSession(msg, phoneNumber, session)
		return
	}

	s.updateSessionActivity(phoneNumber)

	if strings.Contains(strings.ToLower(text), "sayang") {
		res, err := s.genaiSvc.Chat(s.ctx, []string{text})
		if err != nil {
			s.sendReply(msg, "Sorry, I couldn't process your message right now.")
			return
		}

		s.sendReply(msg, res)
	}
}

func (s *WhatsAppBot) handleEndSession(msg *events.Message, phoneNumber string, session *Session) {
	s.sessionsMux.Lock()
	session.WaitingForRating = true
	s.sessionsMux.Unlock()

	s.sendReply(msg, "Terima kasih telah menggunakan layanan kami! 🙏\n\nSilakan berikan rating Anda (1-5):")
}

func (s *WhatsAppBot) handleRatingInput(msg *events.Message, phoneNumber string, text string, session *Session) {
	rating, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil || rating < 1 || rating > 5 {
		s.sendReply(msg, "Rating tidak valid. Silakan masukkan angka antara 1-5:")
		return
	}

	s.sessionsMux.Lock()
	session.Rating = rating
	session.WaitingForRating = false
	session.WaitingForComment = true
	s.sessionsMux.Unlock()

	s.sendReply(msg, fmt.Sprintf("Terima kasih! Anda memberikan rating %d ⭐\n\nSilakan berikan komentar atau saran Anda (atau ketik '/skip' untuk melewati):", rating))
}

func (s *WhatsAppBot) handleCommentInput(msg *events.Message, phoneNumber string, text string, session *Session) {
	ctx := context.Background()

	users, _, err := s.userRepo.List(ctx, &entity.GetUsersFilter{
		Offset: 0,
		Limit:  1,
		Search: phoneNumber,
	})

	if err != nil {
		s.clientLog.Errorf("Failed to find user: %v", err)
		s.sendReply(msg, "Sorry, something went wrong. Please try again later.")
		return
	}

	if len(users) == 0 {
		s.clientLog.Errorf("User not found for phone number: %s", phoneNumber)
		s.sendReply(msg, "Sorry, user not found. Please contact support.")
		return
	}

	user := users[0]

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
		ID:        feedbackID,
		UserID:    user.ID,
		Rating:    session.Rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	if err := s.feedbackRepo.Create(ctx, feedback); err != nil {
		s.clientLog.Errorf("Failed to save feedback: %v", err)
		s.sendReply(msg, "Sorry, something went wrong while saving your feedback.")
		return
	}

	s.deleteSession(phoneNumber)

	s.sendReply(msg, "Terima kasih atas feedback Anda! 🙏\n\nSampai jumpa lagi! 👋")

	log.Info(log.CustomLogInfo{
		"phone_number": phoneNumber,
		"user_id":      user.ID.String(),
		"rating":       session.Rating,
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
