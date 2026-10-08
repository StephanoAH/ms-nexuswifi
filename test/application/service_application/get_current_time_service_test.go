package service_application

import (
	"errors"
	"testing"
	"time"
	"genesis/pos/application/service_application"
	"genesis/pos/domain/entities"
)

// Mock del caso de uso
type mockGetCurrentTimeUseCase struct {
	currentTime time.Time
	err         error
}

func (m *mockGetCurrentTimeUseCase) ExecuteUseCase() (time.Time, error) {
	return m.currentTime, m.err
}

// Mock del caso de uso de ventas (necesario para el servicio)
type mockGetSalesUseCase struct {
	sales string
	err   error
}

func (m *mockGetSalesUseCase) ExecuteUseCase() (string, error) {
	return m.sales, m.err
}

func TestGetCurrentTimeService_ExecuteService_Success(t *testing.T) {
	// Arrange
	expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	mockUseCase := &mockGetCurrentTimeUseCase{
		currentTime: expectedTime,
		err:         nil,
	}
	mockSalesUseCase := &mockGetSalesUseCase{
		sales: "test sales",
		err:   nil,
	}
	service := &service_application.GetCurrentTimeServiceImpl{
		GetCurrentTimeServiceUseCase: mockUseCase,
		GetSalesUseCase:              mockSalesUseCase,
	}

	// Act
	var response entities.ResponseGetCurrentTime
	response = service.ExecuteService()

	// Assert
	if response.StatusCode != 200 {
		t.Errorf("Expected StatusCode 200, got %d", response.StatusCode)
	}
	if response.Message != "The current date is" {
		t.Errorf("Expected message 'The current date is', got '%s'", response.Message)
	}
	if !response.Result.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, response.Result)
	}
}

func TestGetCurrentTimeService_ExecuteService_Error(t *testing.T) {
	// Arrange
	expectedError := errors.New("database error")
	mockUseCase := &mockGetCurrentTimeUseCase{
		currentTime: time.Time{},
		err:         expectedError,
	}
	mockSalesUseCase := &mockGetSalesUseCase{
		sales: "",
		err:   nil,
	}
	service := &service_application.GetCurrentTimeServiceImpl{
		GetCurrentTimeServiceUseCase: mockUseCase,
		GetSalesUseCase:              mockSalesUseCase,
	}

	// Act
	var response entities.ResponseGetCurrentTime
	response = service.ExecuteService()

	// Assert
	if response.StatusCode != 500 {
		t.Errorf("Expected StatusCode 500, got %d", response.StatusCode)
	}
	if response.Message != "Error querying the current date" {
		t.Errorf("Expected message 'Error querying the current date', got '%s'", response.Message)
	}
}

