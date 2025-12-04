package repository

import (
	"QnA/internal/domain"
	"context"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, q *domain.Question) error
	GetAllQuestions(ctx context.Context) ([]domain.Question, error)
	GetQuestionByID(ctx context.Context, id string) (*domain.Question, []domain.Answer, error)
	DeleteQuestion(ctx context.Context, id string) error
}

type questionRepo struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepo{db}
}

func (r *questionRepo) CreateQuestion(ctx context.Context, q *domain.Question) error {
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *questionRepo) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	var q []domain.Question
	err := r.db.WithContext(ctx).Find(&q).Error
	return q, err
}

func (r *questionRepo) GetQuestionByID(ctx context.Context, id string) (*domain.Question, []domain.Answer, error) {
	var question domain.Question
	var answers []domain.Answer

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&question).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.WithContext(ctx).
		Where("question_id = ?", id).
		Find(&answers).Error

	return &question, answers, err
}

func (r *questionRepo) DeleteQuestion(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Question{}).Error
}
