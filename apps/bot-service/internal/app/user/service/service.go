package service

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
)

type UserService struct {
	userRepo  contracts.UserRepository
	validator validator.CustomValidatorInterface
	uuidPkg   uuid.UUIDInterface
}

func NewUserService(userRepo contracts.UserRepository, validatorService validator.CustomValidatorInterface, uuidService uuid.UUIDInterface) *UserService {
	return &UserService{
		userRepo:  userRepo,
		validator: validatorService,
		uuidPkg:   uuidService,
	}
}
