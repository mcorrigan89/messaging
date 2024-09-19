package entities

import (
	"github.com/google/uuid"
)

type Email struct {
	ID        uuid.UUID
	MessageID string
}
