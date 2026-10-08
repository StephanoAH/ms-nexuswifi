package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"genesis/pos/domain/entities"
	"genesis/pos/presentation/container"
	"genesis/pos/presentation/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock del servicio
type mockGetCurrentTimeService struct {
	response entities.ResponseGetCurrentTime
}

func (m *mockGetCurrentTimeService) ExecuteService() entities.ResponseGetCurrentTime {
	return m.response
}

func TestGetCurrentTime_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	expectedResponse := entities.ResponseGetCurrentTime{
		StatusCode: 200,
		Message:    "The current date is",
		Result:     expectedTime,
	}

	// Guardar servicio original
	originalService := container.GetCurrentTimeService

	// Mock del servicio en el container
	container.GetCurrentTimeService = &mockGetCurrentTimeService{
		response: expectedResponse,
	}
	defer func() {
		container.GetCurrentTimeService = originalService
	}()

	// Crear request
	req, _ := http.NewRequest("GET", "/api/current-time", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handlers.GetCurrentTime(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response entities.ResponseHttp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expectedTime.String(), response.Message)
}

func TestGetCurrentTime_Error(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	errorResponse := entities.ResponseGetCurrentTime{
		StatusCode: 500,
		Message:    "Error querying the current date",
		Result:     time.Time{},
	}

	// Guardar servicio original
	originalService := container.GetCurrentTimeService

	// Mock del servicio en el container
	container.GetCurrentTimeService = &mockGetCurrentTimeService{
		response: errorResponse,
	}
	defer func() {
		container.GetCurrentTimeService = originalService
	}()

	// Crear request
	req, _ := http.NewRequest("GET", "/api/current-time", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handlers.GetCurrentTime(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response entities.ResponseHttp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, 500, response.Code)
}

