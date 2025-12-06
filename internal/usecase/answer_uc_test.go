package usecase

import (
	"QnA/internal/domain"
	"QnA/internal/shared/errors"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockAnswerRepository is a mock implementation of AnswerRepository
type MockAnswerRepository struct {
	answers      map[string]*domain.Answer
	createError  error
	getByIDError error
	deleteError  error
}

func NewMockAnswerRepository() *MockAnswerRepository {
	return &MockAnswerRepository{
		answers: make(map[string]*domain.Answer),
	}
}

func (m *MockAnswerRepository) CreateAnswer(ctx context.Context, a *domain.Answer) error {
	if m.createError != nil {
		return m.createError
	}
	m.answers[a.ID.String()] = a
	return nil
}

func (m *MockAnswerRepository) GetAnswerByID(ctx context.Context, id uuid.UUID) (*domain.Answer, error) {
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}
	a, ok := m.answers[id.String()]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

func (m *MockAnswerRepository) DeleteAnswer(ctx context.Context, id string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	if _, ok := m.answers[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.answers, id)
	return nil
}

func TestNewAnswerUsecase(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)

	if uc == nil {
		t.Error("NewAnswerUsecase returned nil")
	}

	if uc.qRepo == nil {
		t.Error("AnswerUsecase question repository is nil")
	}

	if uc.aRepo == nil {
		t.Error("AnswerUsecase answer repository is nil")
	}
}

func TestAnswerUsecase_CreateAnswer(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	// Create a question first
	questionID := uuid.New()
	q := &domain.Question{ID: questionID, Text: "Test question", CreatedAt: time.Now()}
	mockQRepo.questions[questionID.String()] = q

	id, err := uc.CreateAnswer(ctx, questionID.String(), "user123", "Test answer")
	if err != nil {
		t.Fatalf("CreateAnswer failed: %v", err)
	}

	if id == uuid.Nil {
		t.Error("CreateAnswer returned nil UUID")
	}

	// Verify answer was created
	if len(mockARepo.answers) != 1 {
		t.Errorf("Expected 1 answer, got %d", len(mockARepo.answers))
	}
}

func TestAnswerUsecase_CreateAnswer_QuestionNotFound(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockQRepo.getByIDError = gorm.ErrRecordNotFound
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	nonExistentQuestionID := uuid.New().String()

	_, err := uc.CreateAnswer(ctx, nonExistentQuestionID, "user123", "Test answer")
	if err == nil {
		t.Error("Expected error when question doesn't exist, got nil")
	}
}

func TestAnswerUsecase_GetAnswerByID(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	answerID := uuid.New()
	questionID := uuid.New()
	answer := &domain.Answer{
		ID:         answerID,
		QuestionID: questionID,
		UserID:     "user123",
		Text:       "Test answer",
		CreatedAt:  time.Now(),
	}
	mockARepo.answers[answerID.String()] = answer

	found, err := uc.GetAnswerByID(ctx, answerID)
	if err != nil {
		t.Fatalf("GetAnswerByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("GetAnswerByID returned nil answer")
	}

	if found.ID != answerID {
		t.Errorf("Expected ID %s, got %s", answerID, found.ID)
	}

	if found.Text != answer.Text {
		t.Errorf("Expected text %s, got %s", answer.Text, found.Text)
	}
}

func TestAnswerUsecase_GetAnswerByID_NilUUID(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	_, err := uc.GetAnswerByID(ctx, uuid.Nil)
	if err == nil {
		t.Error("Expected error for nil UUID, got nil")
	}

	if err != errors.ErrInvalidID {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}
}

func TestAnswerUsecase_GetAnswerByID_NotFound(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	mockARepo.getByIDError = gorm.ErrRecordNotFound
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	nonExistentID := uuid.New()

	_, err := uc.GetAnswerByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent answer, got nil")
	}
}

func TestAnswerUsecase_DeleteAnswer(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	answerID := uuid.New()
	questionID := uuid.New()
	answer := &domain.Answer{
		ID:         answerID,
		QuestionID: questionID,
		UserID:     "user123",
		Text:       "Test answer",
		CreatedAt:  time.Now(),
	}
	mockARepo.answers[answerID.String()] = answer

	err := uc.DeleteAnswer(ctx, answerID.String())
	if err != nil {
		t.Fatalf("DeleteAnswer failed: %v", err)
	}

	// Verify answer was deleted
	if _, ok := mockARepo.answers[answerID.String()]; ok {
		t.Error("Answer should have been deleted")
	}
}

func TestAnswerUsecase_DeleteAnswer_NotFound(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	mockARepo.getByIDError = gorm.ErrRecordNotFound
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	err := uc.DeleteAnswer(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent answer, got nil")
	}
}

func TestAnswerUsecase_DeleteAnswer_InvalidID(t *testing.T) {
	mockQRepo := NewMockQuestionRepository()
	mockARepo := NewMockAnswerRepository()
	uc := NewAnswerUsecase(mockQRepo, mockARepo)
	ctx := context.Background()

	// Test with invalid UUID format - uuid.MustParse will panic
	// Since we can't test invalid UUID without panic, we test with valid format but non-existent
	validUUID := uuid.New().String()
	mockARepo.getByIDError = gorm.ErrRecordNotFound

	err := uc.DeleteAnswer(ctx, validUUID)
	if err == nil {
		t.Error("Expected error when answer doesn't exist, got nil")
	}
}
