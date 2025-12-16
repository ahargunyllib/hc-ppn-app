package whatsapp

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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

		session = s.createSession(phoneNumber, &chatJID, &userRes.User)
		greeting := getTimeBasedGreeting(getJakartaTime())
		personalizedGreeting := formatUserGreeting(&userRes.User, greeting)
		welcomeMessage := fmt.Sprintf("%s 👋\n\nSelamat datang di *Layanan WhatsApp HC PPN*\n\nSaya adalah asisten virtual yang siap membantu Anda dengan pertanyaan seputar layanan kami.\n\nAda yang bisa saya bantu hari ini?", personalizedGreeting)
		s.sendMessage(chatJID, welcomeMessage)
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

	if strings.ToLower(strings.TrimSpace(text)) == "/help" {
		s.handleHelpCommand(msg)
		return
	}

	if strings.ToLower(strings.TrimSpace(text)) == "/selesai" {
		s.handleEndSession(msg, session)
		return
	}

	// Rate limiting: prevent spam
	s.sessionsMux.Lock()

	// 1. Rapid message detection: block messages < 3 seconds apart
	if len(session.MessageHistory) > 0 {
		lastMsgTime := session.MessageHistory[len(session.MessageHistory)-1]
		timeSinceLastMsg := time.Since(lastMsgTime)
		if timeSinceLastMsg < 3*time.Second {
			s.sessionsMux.Unlock()
			s.sendReply(msg, "Mohon tunggu sebentar sebelum mengirim pesan berikutnya 🙏")
			return
		}
	}

	// 2. Sliding window counter: max 20 messages per 10 minutes
	const maxMessagesInWindow = 20
	const windowDuration = 10 * time.Minute

	now := time.Now()
	session.MessageHistory = filterRecentMessages(session.MessageHistory, now, windowDuration)

	if len(session.MessageHistory) >= maxMessagesInWindow {
		s.sessionsMux.Unlock()
		s.sendReply(msg, "Anda telah mencapai batas maksimal pesan (20 pesan per 10 menit). Mohon tunggu beberapa saat 🙏")
		return
	}

	// Add current message to history
	session.MessageHistory = append(session.MessageHistory, now)

	// Reset feedback prompt if user continues conversation
	// This prevents auto-close when user is actively engaging
	if session.FeedbackPromptSent {
		session.FeedbackPromptSent = false
		session.FeedbackPromptSentAt = nil
		session.IsAutoPrompt = false
	}
	s.sessionsMux.Unlock()

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

	if difyResp.ConversationID != "" {
		s.sessionsMux.Lock()
		if session.ConversationID == "" {
			session.ConversationID = difyResp.ConversationID
		}
		s.sessionsMux.Unlock()
	}

	s.sendReply(msg, difyResp.Answer)
}

func (s *WhatsAppBot) handleEndSession(msg *events.Message, session *Session) {
	s.sessionsMux.Lock()
	session.WaitingForRating = true
	s.sessionsMux.Unlock()

	s.sendReply(msg, "*[Langkah 1/2]* ⭐\n\nTerima kasih telah menggunakan layanan kami! 🙏\n\nSilakan berikan rating Anda (1-5):\n\n*Skala Penilaian:*\n1 = Sangat Tidak Memuaskan\n2 = Tidak Memuaskan\n3 = Cukup Memuaskan\n4 = Memuaskan\n5 = Sangat Memuaskan")
}

func (s *WhatsAppBot) handleRatingInput(msg *events.Message, text string, session *Session) {
	rating, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil || rating < 1 || rating > 5 {
		s.sendReply(msg, "Mohon maaf, rating harus berupa angka dari 1 sampai 5 ya 😊\n\n*Contoh:* ketik angka *3* untuk rating 3 bintang\n\nSilakan coba lagi:")
		return
	}

	s.sessionsMux.Lock()
	session.Rating = rating
	session.WaitingForRating = false
	session.WaitingForComment = true
	s.sessionsMux.Unlock()

	s.sendReply(msg, getRatingConfirmationMessage(rating))
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
		s.sendReply(msg, "Maaf, terjadi kesalahan saat menyimpan feedback Anda. Silakan coba lagi nanti.")
		return
	}

	s.deleteSession(phoneNumber)

	s.sendReply(msg, getGoodbyeMessage(session.Rating, comment != nil))

	log.Info(log.CustomLogInfo{
		"phone_number": phoneNumber,
		"user_id":      userRes.User.ID,
		"rating":       session.Rating,
		"has_comment":  comment != nil,
		"feedback_id":  feedbackRes.ID,
	}, "[WhatsAppBot] Feedback received and saved")
}

// simulateTyping sends a typing indicator and waits for a realistic delay based on message length
func (s *WhatsAppBot) simulateTyping(chatJID types.JID, messageText string) {
	// Send typing presence
	err := s.client.SendChatPresence(s.ctx, chatJID, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	if err != nil {
		s.clientLog.Warnf("Failed to send typing presence: %v", err)
	}

	// Calculate realistic typing delay based on message length
	// Average human typing speed: 40 WPM ≈ 200 characters per minute ≈ 300ms per character
	// We'll use a faster rate (150ms per char) + base delay to feel responsive but human
	const baseDelayMs = 800
	const msPerChar = 15

	messageLength := len([]rune(messageText))
	calculatedDelay := baseDelayMs + (messageLength * msPerChar)

	// Cap the maximum delay to avoid long waits (max 4 seconds)
	const maxDelayMs = 4000
	delay := math.Min(float64(calculatedDelay), maxDelayMs)

	// Minimum delay of 500ms to ensure typing indicator is visible
	const minDelayMs = 500
	delay = math.Max(delay, minDelayMs)

	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Stop typing presence
	err = s.client.SendChatPresence(s.ctx, chatJID, types.ChatPresencePaused, types.ChatPresenceMediaText)
	if err != nil {
		s.clientLog.Warnf("Failed to send paused presence: %v", err)
	}
}

func (s *WhatsAppBot) sendReply(msg *events.Message, text string) {
	// Simulate typing before sending the reply
	s.simulateTyping(msg.Info.Chat, text)

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
	// Simulate typing before sending the message
	s.simulateTyping(to, text)

	_, err := s.client.SendMessage(s.ctx, to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
		},
	})
	if err != nil {
		s.clientLog.Errorf("Failed to send WhatsApp message: " + err.Error())
	}
}

func getRatingConfirmationMessage(rating int) string {
	var emoji, response string

	switch rating {
	case 5:
		emoji = "🌟"
		response = "Terima kasih atas rating sempurna!"
	case 4:
		emoji = "⭐"
		response = "Terima kasih atas rating yang baik!"
	case 3:
		emoji = "⭐"
		response = "Terima kasih atas rating Anda"
	case 2:
		emoji = "💭"
		response = "Terima kasih atas masukan Anda"
	case 1:
		emoji = "💬"
		response = "Terima kasih telah berbagi pengalaman Anda"
	default:
		emoji = "⭐"
		response = "Terima kasih atas rating Anda"
	}

	return fmt.Sprintf("%s %s\nRating: %d/5\n\n*[Langkah 2/2]* 📝\n\nBantu kami lebih baik lagi dengan memberikan komentar atau saran Anda.\n\n💡 Ketik '/skip' jika ingin melewati.", emoji, response, rating)
}

func getGoodbyeMessage(rating int, hasComment bool) string {
	var message string

	if rating >= 4 {
		message = "Senang mendengar pengalaman Anda positif! 😊\n\n"
	} else if rating == 3 {
		message = "Terima kasih atas masukan Anda. Kami akan terus berusaha lebih baik! 💪\n\n"
	} else {
		message = "Mohon maaf atas ketidaknyamanannya. Kami akan segera memperbaiki layanan kami. 🙏\n\n"
	}

	if hasComment {
		message += "Feedback Anda sangat berharga bagi kami dan akan kami gunakan untuk meningkatkan kualitas layanan.\n\n"
	}

	greeting := getTimeBasedGreeting(getJakartaTime())
	message += fmt.Sprintf("Sampai jumpa lagi! 👋\n\n%s dan semoga harimu menyenangkan! ✨", greeting)

	return message
}

func (s *WhatsAppBot) handleHelpCommand(msg *events.Message) {
	helpMessage := "📖 *Panduan Penggunaan Bot*\n\nSaya adalah asisten virtual yang siap membantu Anda 🤖\n\n*Command yang tersedia:*\n• /help - Menampilkan panduan ini\n• /selesai - Mengakhiri sesi dan memberikan feedback\n\nAnda bisa mengirim pertanyaan kapan saja, dan saya akan membantu menjawabnya! 💬"
	s.sendReply(msg, helpMessage)
}

func (s *WhatsAppBot) autoSubmitFeedback(phoneNumber string, chatJID types.JID) {
	ctx := context.Background()

	userRes, err := s.userSvc.GetByPhoneNumber(ctx, &dto.GetUserByPhoneNumberParam{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		log.Debug(log.CustomLogInfo{
			"phone_number": phoneNumber,
			"error":        err.Error(),
		}, "[WhatsAppBot] Failed to get user for auto-feedback submission")
		return
	}

	rating := 5
	_, err = s.feedbackSvc.Create(ctx, &dto.CreateFeedbackRequest{
		UserID:  userRes.User.ID,
		Rating:  rating,
		Comment: nil,
	})
	if err != nil {
		s.clientLog.Errorf("Failed to auto-submit feedback: %v", err)
		return
	}

	confirmationMessage := "Terima kasih! ✨\n\nKarena tidak ada respons, kami mencatat feedback Anda dengan rating 5 bintang ⭐\n\nKami menghargai waktu Anda dan berharap layanan kami memuaskan. Sampai jumpa lagi! 👋"
	s.sendMessage(chatJID, confirmationMessage)

	log.Info(log.CustomLogInfo{
		"phone_number": phoneNumber,
		"user_id":      userRes.User.ID,
		"rating":       rating,
		"auto_submit":  true,
	}, "[WhatsAppBot] Auto-submitted feedback with rating 5")
}
