package genai

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"google.golang.org/genai"
)

var (
	ModelGemini25Flash = "gemini-2.5-flash"
)

type CustomGenAIInterface interface {
	Chat(ctx context.Context, texts []string) (string, error)
}

type CustomGenAIStruct struct {
	client *genai.Client
}

func getGenAI() CustomGenAIInterface {
	client, err := genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey:  env.AppEnv.GoogleAPIKey,
			Backend: genai.BackendGeminiAPI,
		},
	)
	if err != nil {
		log.Fatal(log.CustomLogInfo{
			"error": err.Error(),
		}, "[GenAI][getGenAI] failed to create GenAI client")
	}

	return &CustomGenAIStruct{
		client: client,
	}
}

var GenAI = getGenAI()

func (o *CustomGenAIStruct) Chat(ctx context.Context, texts []string) (string, error) {
	parts := []*genai.Part{}
	for _, text := range texts {
		parts = append(parts, &genai.Part{
			Text: text,
		})
	}

	contents := []*genai.Content{{Parts: parts}}

	res, err := o.client.Models.GenerateContent(
		ctx,
		ModelGemini25Flash,
		contents,
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{
					{Text: "semua keluaran buat whatsapp jadi plain text, jangan pake markdown atau formatting sama sekali"},
				},
			},
		},
	)
	if err != nil {
		log.Error(log.CustomLogInfo{
			"error": err.Error(),
		}, "[GenAI][Chat] failed to generate content")
		return "", err
	}

	log.Debug(log.CustomLogInfo{
		"response": res,
	}, "[GenAI][Chat] generated content successfully")

	return res.Text(), nil
}
