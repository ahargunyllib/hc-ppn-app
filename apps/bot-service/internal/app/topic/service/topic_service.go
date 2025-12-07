package service

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
)

func (s *TopicService) BulkCreate(ctx context.Context, req *dto.BulkCreateTopicsRequest) error {
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	topics := make([]entity.Topic, 0, len(req.Topics))
	for _, item := range req.Topics {
		topics = append(topics, entity.Topic{
			Title: item.Title,
			Count: item.Count,
		})
	}

	if err := s.topicRepo.BulkCreate(ctx, topics); err != nil {
		return err
	}

	return nil
}

func (s *TopicService) GetHotTopics(ctx context.Context) (*dto.GetHotTopicsResponse, error) {
	hotTopics, err := s.topicRepo.GetHotTopics(ctx)
	if err != nil {
		return nil, err
	}

	hotTopicsData := make([]dto.TopicResponse, 0, len(hotTopics))
	for _, topic := range hotTopics {
		hotTopicsData = append(hotTopicsData, dto.ToTopicResponse(&topic))
	}

	res := &dto.GetHotTopicsResponse{
		Topics: hotTopicsData,
	}

	return res, nil
}
