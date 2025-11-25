package ulid

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/oklog/ulid/v2"
)

type CustomULIDInterface interface {
	New() (ulid.ULID, error)
}

type CustomULIDStruct struct {
	mu      sync.Mutex
	entropy *rand.Rand
}

var ULID = getULID()

func getULID() CustomULIDInterface {
	entropy := rand.New(rand.NewSource(int64(time.Now().UnixNano())))

	return &CustomULIDStruct{
		entropy: entropy,
	}
}

func (u *CustomULIDStruct) New() (ulid.ULID, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, u.entropy)
	if err != nil {
		log.Error(log.CustomLogInfo{
			"error": err.Error(),
		}, "[ULID][New] Failed to generate ULID")

		return ulid.ULID{}, err
	}

	return id, nil
}
