package service

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
)

type TopicService struct {
	topicRepo contracts.TopicRepository
	validator validator.CustomValidatorInterface
}

func NewTopicService(topicRepo contracts.TopicRepository, validatorService validator.CustomValidatorInterface) *TopicService {
	return &TopicService{
		topicRepo: topicRepo,
		validator: validatorService,
	}
}
