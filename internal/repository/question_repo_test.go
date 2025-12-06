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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(&domain.Question{}, &domain.Answer{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestNewQuestionRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)

	if repo == nil {
		t.Error("NewQuestionRepository returned nil")
	}
}

func TestQuestionRepository_CreateQuestion(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)
	ctx := context.Background()

	question := &domain.Question{
		ID:        uuid.New(),
		Text:      "Test question",
		CreatedAt: time.Now(),
	}

	err := repo.CreateQuestion(ctx, question)
	if err != nil {
		t.Fatalf("CreateQuestion failed: %v", err)
	}

	// Verify question was created
	var found domain.Question
	err = db.First(&found, "id = ?", question.ID).Error
	if err != nil {
		t.Fatalf("Question not found in database: %v", err)
	}

	if found.Text != question.Text {
		t.Errorf("Expected text %s, got %s", question.Text, found.Text)
	}
}

func TestQuestionRepository_GetAllQuestions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)
	ctx := context.Background()

	// Create test questions
	questions := []domain.Question{
		{ID: uuid.New(), Text: "Question 1", CreatedAt: time.Now()},
		{ID: uuid.New(), Text: "Question 2", CreatedAt: time.Now()},
		{ID: uuid.New(), Text: "Question 3", CreatedAt: time.Now()},
	}

	for _, q := range questions {
		db.Create(&q)
	}

	result, err := repo.GetAllQuestions(ctx)
	if err != nil {
		t.Fatalf("GetAllQuestions failed: %v", err)
	}

	if len(result) != len(questions) {
		t.Errorf("Expected %d questions, got %d", len(questions), len(result))
	}
}

func TestQuestionRepository_GetQuestionByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)
	ctx := context.Background()

	question := domain.Question{
		ID:        uuid.New(),
		Text:      "Test question",
		CreatedAt: time.Now(),
	}
	db.Create(&question)

	// Create some answers for this question
	answers := []domain.Answer{
		{ID: uuid.New(), QuestionID: question.ID, UserID: "user1", Text: "Answer 1", CreatedAt: time.Now()},
		{ID: uuid.New(), QuestionID: question.ID, UserID: "user2", Text: "Answer 2", CreatedAt: time.Now()},
	}
	for _, a := range answers {
		db.Create(&a)
	}

	found, foundAnswers, err := repo.GetQuestionByID(ctx, question.ID.String())
	if err != nil {
		t.Fatalf("GetQuestionByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("GetQuestionByID returned nil question")
	}

	if found.ID != question.ID {
		t.Errorf("Expected ID %s, got %s", question.ID, found.ID)
	}

	if found.Text != question.Text {
		t.Errorf("Expected text %s, got %s", question.Text, found.Text)
	}

	if len(foundAnswers) != len(answers) {
		t.Errorf("Expected %d answers, got %d", len(answers), len(foundAnswers))
	}
}

func TestQuestionRepository_GetQuestionByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	_, _, err := repo.GetQuestionByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent question, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestQuestionRepository_DeleteQuestion(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)
	ctx := context.Background()

	question := domain.Question{
		ID:        uuid.New(),
		Text:      "Test question",
		CreatedAt: time.Now(),
	}
	db.Create(&question)

	err := repo.DeleteQuestion(ctx, question.ID.String())
	if err != nil {
		t.Fatalf("DeleteQuestion failed: %v", err)
	}

	// Verify question was deleted
	var found domain.Question
	err = db.First(&found, "id = ?", question.ID).Error
	if err != gorm.ErrRecordNotFound {
		t.Error("Question should have been deleted")
	}
}

func TestQuestionRepository_DeleteQuestion_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewQuestionRepository(db)
	ctx := context.Background()

	nonExistentID := uuid.New().String()

	err := repo.DeleteQuestion(ctx, nonExistentID)
	// GORM doesn't return error if record doesn't exist when deleting
	// So this should not error
	if err != nil {
		t.Errorf("DeleteQuestion should not error for non-existent question, got %v", err)
	}
}
