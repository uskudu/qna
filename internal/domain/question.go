package domain

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
