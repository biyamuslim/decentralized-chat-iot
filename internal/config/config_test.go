package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// Test LoadEnv function
func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file for testing
	envFile := ".env.test"
	envContent := "TEST_KEY=HelloWorld"

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(envFile) // Clean up after test

	// Explicitly load the test .env file
	err = godotenv.Load(envFile)
	if err != nil {
		t.Fatalf("Error loading test .env file: %v", err)
	}

	// Retrieve the environment variable
	value := GetEnv("TEST_KEY")

	// Check if the variable is loaded correctly
	if value != "HelloWorld" {
		t.Errorf("Expected 'HelloWorld', but got '%s'", value)
	}
}

// Test GetEnv function
func TestGetEnv(t *testing.T) {
	// Set an environment variable manually
	os.Setenv("SAMPLE_KEY", "SampleValue")

	// Get the value using GetEnv function
	value := GetEnv("SAMPLE_KEY")

	// Verify it returns the correct value
	if value != "SampleValue" {
		t.Errorf("Expected 'SampleValue', but got '%s'", value)
	}
}
