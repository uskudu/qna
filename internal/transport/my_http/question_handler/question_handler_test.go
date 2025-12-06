package question_handler

import (
	"QnA/internal/domain"
	"QnA/internal/usecase"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockQuestionRepositoryForHandler is a mock repository for handler tests
type MockQuestionRepositoryForHandler struct {
	questions    map[string]*domain.Question
	answers      map[string][]domain.Answer
	createError  error
	getAllError  error
	getByIDError error
	deleteError  error
}

func NewMockQuestionRepositoryForHandler() *MockQuestionRepositoryForHandler {
	return &MockQuestionRepositoryForHandler{
		questions: make(map[string]*domain.Question),
		answers:   make(map[string][]domain.Answer),
	}
}

func (m *MockQuestionRepositoryForHandler) CreateQuestion(ctx context.Context, q *domain.Question) error {
	if m.createError != nil {
		return m.createError
	}
	m.questions[q.ID.String()] = q
	return nil
}

func (m *MockQuestionRepositoryForHandler) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}
	result := make([]domain.Question, 0, len(m.questions))
	for _, q := range m.questions {
		result = append(result, *q)
	}
	return result, nil
}

func (m *MockQuestionRepositoryForHandler) GetQuestionByID(ctx context.Context, id string) (*domain.Question, []domain.Answer, error) {
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

func (m *MockQuestionRepositoryForHandler) DeleteQuestion(ctx context.Context, id string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	if _, ok := m.questions[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.questions, id)
	return nil
}

func TestNewQuestionHandler(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	if handler == nil {
		t.Error("NewQuestionHandler returned nil")
	}

	if handler.uc == nil {
		t.Error("QuestionHandler usecase is nil")
	}
}

func TestQuestionHandler_Create(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	reqBody := map[string]string{"text": "What is Go?"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["id"] == nil {
		t.Error("Response missing 'id' field")
	}

	// Verify question was created
	if len(mockRepo.questions) != 1 {
		t.Errorf("Expected 1 question in repository, got %d", len(mockRepo.questions))
	}
}

func TestQuestionHandler_Create_InvalidMethod(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestQuestionHandler_Create_InvalidBody(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestQuestionHandler_Create_EmptyText(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	reqBody := map[string]string{"text": ""}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestQuestionHandler_Create_UsecaseError(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	mockRepo.createError = gorm.ErrInvalidDB
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	reqBody := map[string]string{"text": "What is Go?"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestQuestionHandler_GetAll(t *testing.T) {
	expectedQuestions := []*domain.Question{
		{ID: uuid.New(), Text: "Question 1", CreatedAt: time.Now()},
		{ID: uuid.New(), Text: "Question 2", CreatedAt: time.Now()},
	}
	mockRepo := NewMockQuestionRepositoryForHandler()
	for _, q := range expectedQuestions {
		mockRepo.questions[q.ID.String()] = q
	}
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []domain.Question
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != len(expectedQuestions) {
		t.Errorf("Expected %d questions, got %d", len(expectedQuestions), len(response))
	}
}

func TestQuestionHandler_GetAll_InvalidMethod(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/questions", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestQuestionHandler_GetAll_UsecaseError(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	mockRepo.getAllError = gorm.ErrInvalidDB
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestQuestionHandler_GetByID(t *testing.T) {
	questionID := uuid.New()
	expectedQuestion := &domain.Question{
		ID:        questionID,
		Text:      "Test question",
		CreatedAt: time.Now(),
	}
	expectedAnswers := []domain.Answer{
		{ID: uuid.New(), QuestionID: questionID, UserID: "user1", Text: "Answer 1", CreatedAt: time.Now()},
	}
	mockRepo := NewMockQuestionRepositoryForHandler()
	mockRepo.questions[questionID.String()] = expectedQuestion
	mockRepo.answers[questionID.String()] = expectedAnswers
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/questions/"+questionID.String(), nil)
	req.SetPathValue("id", questionID.String())
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["question"] == nil {
		t.Error("Response missing 'question' field")
	}

	if response["answers"] == nil {
		t.Error("Response missing 'answers' field")
	}
}

func TestQuestionHandler_GetByID_InvalidMethod(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodPost, "/questions/123", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestQuestionHandler_GetByID_NotFound(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	mockRepo.getByIDError = gorm.ErrRecordNotFound
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	nonExistentID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/questions/"+nonExistentID, nil)
	req.SetPathValue("id", nonExistentID)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestQuestionHandler_Delete(t *testing.T) {
	questionID := uuid.New()
	q := &domain.Question{ID: questionID, Text: "Test question", CreatedAt: time.Now()}
	mockRepo := NewMockQuestionRepositoryForHandler()
	mockRepo.questions[questionID.String()] = q
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodDelete, "/questions/"+questionID.String(), nil)
	req.SetPathValue("id", questionID.String())
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify question was deleted
	if _, ok := mockRepo.questions[questionID.String()]; ok {
		t.Error("Question should have been deleted")
	}
}

func TestQuestionHandler_Delete_InvalidMethod(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/questions/123", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestQuestionHandler_Delete_NotFound(t *testing.T) {
	mockRepo := NewMockQuestionRepositoryForHandler()
	mockRepo.getByIDError = gorm.ErrRecordNotFound
	uc := usecase.NewQuestionUsecase(mockRepo)
	handler := NewQuestionHandler(uc)

	nonExistentID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/questions/"+nonExistentID, nil)
	req.SetPathValue("id", nonExistentID)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
