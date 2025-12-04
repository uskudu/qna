package usecase

import (
	"QnA/internal/domain"
	"QnA/internal/repository"
	"QnA/internal/shared/errors"
	"context"
	"time"

	"github.com/google/uuid"
)

type QuestionUsecase struct {
	qRepo repository.QuestionRepository
}

func NewQuestionUsecase(r repository.QuestionRepository) *QuestionUsecase {
	return &QuestionUsecase{qRepo: r}
}

func (u *QuestionUsecase) CreateQuestion(ctx context.Context, text string) (uuid.UUID, error) {
	q := domain.Question{
		ID:        uuid.New(),
		Text:      text,
		CreatedAt: time.Now(),
	}
	err := u.qRepo.CreateQuestion(ctx, &q)
	return q.ID, err
}

func (u *QuestionUsecase) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	return u.qRepo.GetAllQuestions(ctx)
}

func (u *QuestionUsecase) GetQuestionByID(ctx context.Context, id string) (*domain.Question, []domain.Answer, error) {
	questionID := uuid.MustParse(id)
	if questionID == uuid.Nil {
		return &domain.Question{}, nil, errors.ErrInvalidID
	}
	return u.qRepo.GetQuestionByID(ctx, id)
}

func (u *QuestionUsecase) DeleteQuestion(ctx context.Context, id string) error {
	questionID := uuid.MustParse(id)
	if questionID == uuid.Nil {
		return errors.ErrInvalidID
	}

	// проверка существования
	_, _, err := u.qRepo.GetQuestionByID(ctx, id)
	if err != nil {
		return err
	}

	return u.qRepo.DeleteQuestion(ctx, id)
}
