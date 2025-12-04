package repository

import (
	"QnA/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerRepository interface {
	CreateAnswer(ctx context.Context, a *domain.Answer) error
	GetAnswerByID(ctx context.Context, id uuid.UUID) (*domain.Answer, error)
	DeleteAnswer(ctx context.Context, id string) error
}

type answerRepo struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepo{db}
}

func (r *answerRepo) CreateAnswer(ctx context.Context, a *domain.Answer) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *answerRepo) GetAnswerByID(ctx context.Context, id uuid.UUID) (*domain.Answer, error) {
	var answer domain.Answer

	err := r.db.WithContext(ctx).First(&answer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &answer, err
}

func (r *answerRepo) DeleteAnswer(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Answer{}).Error
}
