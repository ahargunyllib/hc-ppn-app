package uuid

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/mock_uuid.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid UUIDInterface

type UUIDInterface interface {
	NewV7() (uuid.UUID, error)
	Parse(string) (uuid.UUID, error)
}

type UUIDStruct struct{}

var UUID = getUUID()

func getUUID() UUIDInterface {
	return &UUIDStruct{}
}

func (u *UUIDStruct) NewV7() (uuid.UUID, error) {
	id, err := uuid.NewV7()

	if err != nil {
		log.Error(log.CustomLogInfo{
			"error": err.Error(),
		}, "[UUID][New] failed to create uuid v7")

		return uuid.Nil, err
	}

	return id, err
}

func (u *UUIDStruct) Parse(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)

	if err != nil {
		log.Error(log.CustomLogInfo{
			"error": err.Error(),
		}, "[UUID][Parse] failed to parse uuid from string")

		return uuid.Nil, err
	}

	return id, nil
}
