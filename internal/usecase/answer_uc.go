package usecase

import (
	"QnA/internal/domain"
	"QnA/internal/repository"
	"QnA/internal/shared/errors"
	"context"
	"time"

	"github.com/google/uuid"
)

type AnswerUsecase struct {
	qRepo repository.QuestionRepository
	aRepo repository.AnswerRepository
}

func NewAnswerUsecase(q repository.QuestionRepository, a repository.AnswerRepository) *AnswerUsecase {
	return &AnswerUsecase{qRepo: q, aRepo: a}
}

func (u *AnswerUsecase) CreateAnswer(ctx context.Context, qid, userID, text string) (uuid.UUID, error) {
	_, _, err := u.qRepo.GetQuestionByID(ctx, qid)
	if err != nil {
		return uuid.Nil, err
	}
	questionID := uuid.MustParse(qid)
	a := domain.Answer{
		ID:         uuid.New(),
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
		CreatedAt:  time.Now(),
	}
	err = u.aRepo.CreateAnswer(ctx, &a)
	return a.ID, err
}

func (u *AnswerUsecase) GetAnswerByID(ctx context.Context, id uuid.UUID) (*domain.Answer, error) {
	//answerID := uuid.MustParse(id)
	if id == uuid.Nil {
		return nil, errors.ErrInvalidID
	}
	return u.aRepo.GetAnswerByID(ctx, id)
}

func (u *AnswerUsecase) DeleteAnswer(ctx context.Context, id string) error {
	answerID := uuid.MustParse(id)
	if answerID == uuid.Nil {
		return errors.ErrInvalidID
	}

	// проверка существования
	_, err := u.aRepo.GetAnswerByID(ctx, answerID)
	if err != nil {
		return err
	}

	return u.aRepo.DeleteAnswer(ctx, id)
}
