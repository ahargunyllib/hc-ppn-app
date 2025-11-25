package whatsapp

import (
	"strings"

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

	text := msg.Message.GetConversation()
	phoneNumber := msg.Info.Sender.User

	log.Debug(log.CustomLogInfo{
		"from": phoneNumber,
		"text": text,
		"meta": meta,
	}, "[WhatsAppBot] Received WhatsApp message")

	if strings.Contains(strings.ToLower(text), "sayang") {
		res, err := s.genaiSvc.Chat(s.ctx, []string{text})
		if err != nil {
			s.sendReply(msg, "Sorry, I couldn't process your message right now.")
			return
		}

		s.sendReply(msg, res)
	}
}

func (s *WhatsAppBot) sendReply(msg *events.Message, text string) {
	_, err := s.client.SendMessage(s.ctx, msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    proto.String(msg.Info.ID),
				Participant: proto.String(msg.Info.Sender.String()),
				QuotedMessage: &waE2E.Message{
					Conversation: proto.String(msg.Message.GetConversation()),
				},
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
