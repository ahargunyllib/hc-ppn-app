package service

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
)

type FeedbackService struct {
	feedbackRepo contracts.FeedbackRepository
	validator    validator.CustomValidatorInterface
	uuidPkg      uuid.UUIDInterface
}

func NewFeedbackService(
	feedbackRepo contracts.FeedbackRepository,
	validatorService validator.CustomValidatorInterface,
	uuidService uuid.UUIDInterface,
) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
		validator:    validatorService,
		uuidPkg:      uuidService,
	}
}
