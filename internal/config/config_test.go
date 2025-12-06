package config

import (
	"os"
	"testing"
)

func TestLoad_WithEnvFile(t *testing.T) {
	// Create a temporary .env file
	envContent := `DATABASE_URL=postgres://user:pass@localhost:5432/testdb
HTTP_PORT=3000
`
	err := os.WriteFile(".env.test", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(".env.test")

	// Save original .env if it exists
	if _, err := os.Stat(".env"); err == nil {
		os.Rename(".env", ".env.backup")
		defer os.Rename(".env.backup", ".env")
	}

	// Copy test env to .env
	err = os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	cfg := Load()

	if cfg == nil {
		t.Fatal("Load returned nil")
	}

	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("Expected DatabaseURL 'postgres://user:pass@localhost:5432/testdb', got %s", cfg.DatabaseURL)
	}

	if cfg.HTTPPort != "3000" {
		t.Errorf("Expected HTTPPort '3000', got %s", cfg.HTTPPort)
	}
}

func TestLoad_WithoutEnvFile(t *testing.T) {
	// Save original .env if it exists
	if _, err := os.Stat(".env"); err == nil {
		os.Rename(".env", ".env.backup")
		defer os.Rename(".env.backup", ".env")
	}

	// Set environment variables
	os.Setenv("DATABASE_URL", "postgres://env:pass@localhost:5432/envdb")
	os.Setenv("HTTP_PORT", "8080")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("HTTP_PORT")
	}()

	cfg := Load()

	if cfg == nil {
		t.Fatal("Load returned nil")
	}

	if cfg.DatabaseURL != "postgres://env:pass@localhost:5432/envdb" {
		t.Errorf("Expected DatabaseURL 'postgres://env:pass@localhost:5432/envdb', got %s", cfg.DatabaseURL)
	}

	if cfg.HTTPPort != "8080" {
		t.Errorf("Expected HTTPPort '8080', got %s", cfg.HTTPPort)
	}
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Save original .env if it exists
	if _, err := os.Stat(".env"); err == nil {
		os.Rename(".env", ".env.backup")
		defer os.Rename(".env.backup", ".env")
	}

	// Set environment variables
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("HTTP_PORT", "9000")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("HTTP_PORT")
	}()

	cfg := Load()

	if cfg == nil {
		t.Fatal("Load returned nil")
	}

	if cfg.DatabaseURL == "" {
		t.Error("DatabaseURL should not be empty")
	}

	if cfg.HTTPPort == "" {
		t.Error("HTTPPort should not be empty")
	}
}

func TestLoad_EmptyEnvironment(t *testing.T) {
	// Save original .env if it exists
	if _, err := os.Stat(".env"); err == nil {
		os.Rename(".env", ".env.backup")
		defer os.Rename(".env.backup", ".env")
	}

	// Unset environment variables
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("HTTP_PORT")

	cfg := Load()

	if cfg == nil {
		t.Fatal("Load returned nil")
	}

	// Config should still be created, but with empty values
	if cfg.DatabaseURL != "" {
		t.Errorf("Expected empty DatabaseURL, got %s", cfg.DatabaseURL)
	}

	if cfg.HTTPPort != "" {
		t.Errorf("Expected empty HTTPPort, got %s", cfg.HTTPPort)
	}
}
