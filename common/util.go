package common

import (
	"github.com/go-basic/uuid"
)

func GetUuid() string {
	return uuid.New()
}
