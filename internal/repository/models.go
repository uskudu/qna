package repository

import (
	"time"

	"github.com/google/uuid"
)

type QuestionModel struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Text      string
	CreatedAt time.Time
	Answers   []AnswerModel `gorm:"constraint:OnDelete:CASCADE;"`
}

type AnswerModel struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	QuestionID uuid.UUID
	UserID     string
	Text       string
	CreatedAt  time.Time
}
