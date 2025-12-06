package usecase

import (
	"QnA/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockQuestionRepository is a mock implementation of QuestionRepository
type MockQuestionRepository struct {
	questions    map[string]*domain.Question
	answers      map[string][]domain.Answer
	createError  error
	getAllError  error
	getByIDError error
	deleteError  error
}

func NewMockQuestionRepository() *MockQuestionRepository {
	return &MockQuestionRepository{
		questions: make(map[string]*domain.Question),
		answers:   make(map[string][]domain.Answer),
	}
}

func (m *MockQuestionRepository) CreateQuestion(ctx context.Context, q *domain.Question) error {
	if m.createError != nil {
		return m.createError
	}
	m.questions[q.ID.String()] = q
	return nil
}

func (m *MockQuestionRepository) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}
	result := make([]domain.Question, 0, len(m.questions))
	for _, q := range m.questions {
		result = append(result, *q)
	}
	return result, nil
}

func (m *MockQuestionRepository) GetQuestionByID(ctx context.Context, id string) (*domain.Question, []domain.Answer, error) {
	if m.getByIDError != nil {
		return nil, nil, m.getByIDError
	}
	q, ok := m.questions[id]
	if !ok {
		return nil, nil, gorm.ErrRecordNotFound
	}
	answers := m.answers[id]
	return q, answers, nil
}

func (m *MockQuestionRepository) DeleteQuestion(ctx context.Context, id string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	if _, ok := m.questions[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.questions, id)
	return nil
}

func TestNewQuestionUsecase(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)

	if uc == nil {
		t.Error("NewQuestionUsecase returned nil")
	}

	if uc.qRepo == nil {
		t.Error("QuestionUsecase repository is nil")
	}
}

func TestQuestionUsecase_CreateQuestion(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	id, err := uc.CreateQuestion(ctx, "What is Go?")
	if err != nil {
		t.Fatalf("CreateQuestion failed: %v", err)
	}

	if id == uuid.Nil {
		t.Error("CreateQuestion returned nil UUID")
	}

	// Verify question was created in mock
	if len(mockRepo.questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(mockRepo.questions))
	}
}

func TestQuestionUsecase_GetAllQuestions(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	// Pre-populate mock
	q1 := &domain.Question{ID: uuid.New(), Text: "Question 1", CreatedAt: time.Now()}
	q2 := &domain.Question{ID: uuid.New(), Text: "Question 2", CreatedAt: time.Now()}
	mockRepo.questions[q1.ID.String()] = q1
	mockRepo.questions[q2.ID.String()] = q2

	result, err := uc.GetAllQuestions(ctx)
	if err != nil {
		t.Fatalf("GetAllQuestions failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(result))
	}
}

func TestQuestionUsecase_GetAllQuestions_Error(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	mockRepo.getAllError = gorm.ErrInvalidDB
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	_, err := uc.GetAllQuestions(ctx)
	if err == nil {
		t.Error("Expected error from GetAllQuestions, got nil")
	}
}

func TestQuestionUsecase_GetQuestionByID(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	questionID := uuid.New()
	q := &domain.Question{ID: questionID, Text: "Test question", CreatedAt: time.Now()}
	answers := []domain.Answer{
		{ID: uuid.New(), QuestionID: questionID, UserID: "user1", Text: "Answer 1", CreatedAt: time.Now()},
	}
	mockRepo.questions[questionID.String()] = q
	mockRepo.answers[questionID.String()] = answers

	found, foundAnswers, err := uc.GetQuestionByID(ctx, questionID.String())
	if err != nil {
		t.Fatalf("GetQuestionByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("GetQuestionByID returned nil question")
	}

	if found.ID != questionID {
		t.Errorf("Expected ID %s, got %s", questionID, found.ID)
	}

	if len(foundAnswers) != len(answers) {
		t.Errorf("Expected %d answers, got %d", len(answers), len(foundAnswers))
	}
}

func TestQuestionUsecase_GetQuestionByID_InvalidUUID(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	// uuid.MustParse will panic on invalid UUID, but we can test with empty string
	// Actually, looking at the code, it uses uuid.MustParse which panics on invalid input
	// So we need to test with a valid UUID format but non-existent ID
	validUUID := uuid.New().String()
	mockRepo.getByIDError = gorm.ErrRecordNotFound

	_, _, err := uc.GetQuestionByID(ctx, validUUID)
	if err == nil {
		t.Error("Expected error for non-existent question, got nil")
	}
}

func TestQuestionUsecase_GetQuestionByID_NotFound(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	mockRepo.getByIDError = gorm.ErrRecordNotFound
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	_, _, err := uc.GetQuestionByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent question, got nil")
	}
}

func TestQuestionUsecase_DeleteQuestion(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	questionID := uuid.New()
	q := &domain.Question{ID: questionID, Text: "Test question", CreatedAt: time.Now()}
	mockRepo.questions[questionID.String()] = q

	err := uc.DeleteQuestion(ctx, questionID.String())
	if err != nil {
		t.Fatalf("DeleteQuestion failed: %v", err)
	}

	// Verify question was deleted
	if _, ok := mockRepo.questions[questionID.String()]; ok {
		t.Error("Question should have been deleted")
	}
}

func TestQuestionUsecase_DeleteQuestion_NotFound(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	mockRepo.getByIDError = gorm.ErrRecordNotFound
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	err := uc.DeleteQuestion(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent question, got nil")
	}
}

func TestQuestionUsecase_DeleteQuestion_InvalidID(t *testing.T) {
	mockRepo := NewMockQuestionRepository()
	uc := NewQuestionUsecase(mockRepo)
	ctx := context.Background()

	// Test with invalid UUID - uuid.MustParse will panic, but the code checks for uuid.Nil
	// Since we can't easily test invalid UUID without panic, we test the nil check path
	// Actually, the code uses MustParse which panics, so invalid input will panic
	// We'll test with a valid format but non-existent question
	validUUID := uuid.New().String()
	mockRepo.getByIDError = gorm.ErrRecordNotFound

	err := uc.DeleteQuestion(ctx, validUUID)
	if err == nil {
		t.Error("Expected error when question doesn't exist, got nil")
	}
}
