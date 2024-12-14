package internal

import (
	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	generatedUUID, err := uuid.NewUUID()

	if err != nil {
		panic(err)
	}

	return generatedUUID
}
