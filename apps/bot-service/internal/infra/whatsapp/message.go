package whatsapp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/dify"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/phoneutil"
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

	phoneNumber := phoneutil.NormalizeToE164(msg.Info.Sender.User)
	chatJID := msg.Info.Chat

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

	session := s.getSession(phoneNumber)
	if session == nil {
		userRes, err := s.userSvc.GetByPhoneNumber(s.ctx, &dto.GetUserByPhoneNumberParam{
			PhoneNumber: phoneNumber,
		})
		if err != nil {
			// Silently ignore unauthorized phone numbers
			log.Debug(log.CustomLogInfo{
				"phone_number": phoneNumber,
				"error":        err.Error(),
			}, "[WhatsAppBot] Unauthorized phone number attempted to start session")
			return
		}

		log.Info(log.CustomLogInfo{
			"phone_number": phoneNumber,
			"user_id":      userRes.User.ID,
		}, "[WhatsAppBot] Starting new session for authorized phone number")

		session = s.createSession(phoneNumber, &chatJID)
		s.sendMessage(chatJID, "Halo! Selamat datang di layanan WhatsApp kami. Ada yang bisa kami bantu?")
		return
	}

	if session.WaitingForRating {
		s.handleRatingInput(msg, text, session)
		return
	}

	if session.WaitingForComment {
		s.handleCommentInput(msg, phoneNumber, text, session)
		return
	}

	if strings.ToLower(strings.TrimSpace(text)) == "/selesai" {
		s.handleEndSession(msg, session)
		return
	}

	s.updateSessionActivity(phoneNumber)

	difyReq := &dify.Request{
		Inputs:         make(map[string]any),
		Query:          text,
		ResponseMode:   "blocking",
		ConversationID: session.ConversationID,
		User:           phoneNumber,
		Files:          []any{},
	}

	log.Debug(log.CustomLogInfo{
		"difyReq": difyReq,
	}, "[WhatsAppBot] Sending message to Dify AI")

	difyResp, err := s.difySvc.ChatMessages(s.ctx, difyReq)
	if err != nil {
		s.clientLog.Errorf("Failed to get response from Dify AI: %v", err)
		s.sendReply(msg, "Maaf, saya tidak dapat memproses pesan Anda saat ini. Silakan coba lagi nanti.")
		return
	}

	log.Debug(log.CustomLogInfo{
		"difyResp": difyResp,
	}, "[WhatsAppBot] Received response from Dify AI")

	if difyResp.ConversationID != "" && session.ConversationID == "" {
		s.sessionsMux.Lock()
		session.ConversationID = difyResp.ConversationID
		s.sessionsMux.Unlock()
	}

	s.sendReply(msg, difyResp.Answer)
}

func (s *WhatsAppBot) handleEndSession(msg *events.Message, session *Session) {
	s.sessionsMux.Lock()
	session.WaitingForRating = true
	s.sessionsMux.Unlock()

	s.sendReply(msg, "Terima kasih telah menggunakan layanan kami! 🙏\n\nSilakan berikan rating Anda (1-5):")
}

func (s *WhatsAppBot) handleRatingInput(msg *events.Message, text string, session *Session) {
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

	userRes, err := s.userSvc.GetByPhoneNumber(ctx, &dto.GetUserByPhoneNumberParam{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		// Silently ignore unauthorized phone numbers
		log.Debug(log.CustomLogInfo{
			"phone_number": phoneNumber,
			"error":        err.Error(),
		}, "[WhatsAppBot] Unauthorized phone number attempted feedback submission")
		s.deleteSession(phoneNumber)
		return
	}

	var comment *string
	trimmedText := strings.TrimSpace(text)
	if strings.ToLower(trimmedText) != "/skip" && trimmedText != "" {
		comment = &trimmedText
	}

	feedbackRes, err := s.feedbackSvc.Create(ctx, &dto.CreateFeedbackRequest{
		UserID:  userRes.User.ID,
		Rating:  session.Rating,
		Comment: comment,
	})
	if err != nil {
		s.clientLog.Errorf("Failed to save feedback: %v", err)
		s.sendReply(msg, "Sorry, something went wrong while saving your feedback.")
		return
	}

	s.deleteSession(phoneNumber)

	s.sendReply(msg, "Terima kasih atas feedback Anda! 🙏\n\nSampai jumpa lagi! 👋")

	log.Info(log.CustomLogInfo{
		"phone_number": phoneNumber,
		"user_id":      userRes.User.ID,
		"rating":       session.Rating,
		"has_comment":  comment != nil,
		"feedback_id":  feedbackRes.ID,
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
