package repository

import (
	"QnA/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForAnswer(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&domain.Question{}, &domain.Answer{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestNewAnswerRepository(t *testing.T) {
	db := setupTestDBForAnswer(t)
	repo := NewAnswerRepository(db)

	if repo == nil {
		t.Error("NewAnswerRepository returned nil")
	}
}

func TestAnswerRepository_CreateAnswer(t *testing.T) {
	db := setupTestDBForAnswer(t)
	repo := NewAnswerRepository(db)
	ctx := context.Background()

	// Create a question first
	question := domain.Question{
		ID:        uuid.New(),
		Text:      "Test question",
		CreatedAt: time.Now(),
	}
	db.Create(&question)

	answer := &domain.Answer{
		ID:         uuid.New(),
		QuestionID: question.ID,
		UserID:     "user123",
		Text:       "Test answer",
		CreatedAt:  time.Now(),
	}

	err := repo.CreateAnswer(ctx, answer)
	if err != nil {
		t.Fatalf("CreateAnswer failed: %v", err)
	}

	// Verify answer was created
	var found domain.Answer
	err = db.First(&found, "id = ?", answer.ID).Error
	if err != nil {
		t.Fatalf("Answer not found in database: %v", err)
	}

	if found.Text != answer.Text {
		t.Errorf("Expected text %s, got %s", answer.Text, found.Text)
	}

	if found.UserID != answer.UserID {
		t.Errorf("Expected user_id %s, got %s", answer.UserID, found.UserID)
	}

	if found.QuestionID != answer.QuestionID {
		t.Errorf("Expected question_id %s, got %s", answer.QuestionID, found.QuestionID)
	}
}

func TestAnswerRepository_GetAnswerByID(t *testing.T) {
	db := setupTestDBForAnswer(t)
	repo := NewAnswerRepository(db)
	ctx := context.Background()

	// Create a question first
	question := domain.Question{
		ID:        uuid.New(),
		Text:      "Test question",
		CreatedAt: time.Now(),
	}
	db.Create(&question)

	answer := domain.Answer{
		ID:         uuid.New(),
		QuestionID: question.ID,
		UserID:     "user123",
		Text:       "Test answer",
		CreatedAt:  time.Now(),
	}
	db.Create(&answer)

	found, err := repo.GetAnswerByID(ctx, answer.ID)
	if err != nil {
		t.Fatalf("GetAnswerByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("GetAnswerByID returned nil answer")
	}

	if found.ID != answer.ID {
		t.Errorf("Expected ID %s, got %s", answer.ID, found.ID)
	}

	if found.Text != answer.Text {
		t.Errorf("Expected text %s, got %s", answer.Text, found.Text)
	}
}

func TestAnswerRepository_GetAnswerByID_NotFound(t *testing.T) {
	db := setupTestDBForAnswer(t)
	repo := NewAnswerRepository(db)
	ctx := context.Background()

	nonExistentID := uuid.New()

	_, err := repo.GetAnswerByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent answer, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestAnswerRepository_DeleteAnswer(t *testing.T) {
	db := setupTestDBForAnswer(t)
	repo := NewAnswerRepository(db)
	ctx := context.Background()

	// Create a question first
	question := domain.Question{
		ID:        uuid.New(),
		Text:      "Test question",
		CreatedAt: time.Now(),
	}
	db.Create(&question)

	answer := domain.Answer{
		ID:         uuid.New(),
		QuestionID: question.ID,
		UserID:     "user123",
		Text:       "Test answer",
		CreatedAt:  time.Now(),
	}
	db.Create(&answer)

	err := repo.DeleteAnswer(ctx, answer.ID.String())
	if err != nil {
		t.Fatalf("DeleteAnswer failed: %v", err)
	}

	// Verify answer was deleted
	var found domain.Answer
	err = db.First(&found, "id = ?", answer.ID).Error
	if err != gorm.ErrRecordNotFound {
		t.Error("Answer should have been deleted")
	}
}

func TestAnswerRepository_DeleteAnswer_NotFound(t *testing.T) {
	db := setupTestDBForAnswer(t)
	repo := NewAnswerRepository(db)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	err := repo.DeleteAnswer(ctx, nonExistentID)
	// GORM doesn't return error if record doesn't exist when deleting
	if err != nil {
		t.Errorf("DeleteAnswer should not error for non-existent answer, got %v", err)
	}
}
