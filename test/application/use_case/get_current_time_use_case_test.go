package use_case

import (
	"errors"
	"testing"
	"time"
	"genesis/pos/application/use_case"
)

// Mock del repositorio
type mockGetCurrentTimeRepository struct {
	currentTime time.Time
	err         error
}

func (m *mockGetCurrentTimeRepository) ExecuteImpl() (time.Time, error) {
	return m.currentTime, m.err
}

func TestGetCurrentTimeUseCase_ExecuteUseCase_Success(t *testing.T) {
	// Arrange
	expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := &mockGetCurrentTimeRepository{
		currentTime: expectedTime,
		err:         nil,
	}
	useCase := &use_case.GetCurrentTimeUseCaseImpl{
		GetCurrentTimeRepo: mockRepo,
	}

	// Act
	result, err := useCase.ExecuteUseCase()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !result.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, result)
	}
}

func TestGetCurrentTimeUseCase_ExecuteUseCase_Error(t *testing.T) {
	// Arrange
	expectedError := errors.New("database error")
	mockRepo := &mockGetCurrentTimeRepository{
		currentTime: time.Time{},
		err:         expectedError,
	}
	useCase := &use_case.GetCurrentTimeUseCaseImpl{
		GetCurrentTimeRepo: mockRepo,
	}

	// Act
	result, err := useCase.ExecuteUseCase()

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
	if !result.IsZero() {
		t.Errorf("Expected zero time, got %v", result)
	}
}

