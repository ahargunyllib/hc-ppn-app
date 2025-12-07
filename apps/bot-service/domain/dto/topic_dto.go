package dto

import (
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
)

type TopicRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
	Count int    `json:"count" validate:"required,min=1"`
}

type BulkCreateTopicsRequest struct {
	Topics []TopicRequest `json:"topics" validate:"required,min=1,max=100,dive"`
}

type TopicResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Count     int    `json:"count"`
	CreatedAt string `json:"createdAt"`
}

type GetHotTopicsResponse struct {
	Topics []TopicResponse `json:"topics"`
}

type GetTopicsCountResponse struct {
	TotalTopics int `json:"totalTopics"`
}

func ToTopicResponse(topic *entity.Topic) TopicResponse {
	return TopicResponse{
		ID:        topic.ID,
		Title:     topic.Title,
		Count:     topic.Count,
		CreatedAt: topic.CreatedAt.Format(time.RFC3339),
	}
}
