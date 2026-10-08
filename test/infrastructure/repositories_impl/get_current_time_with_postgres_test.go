package repositories_impl

import (
	"errors"
	"testing"
	"genesis/pos/infrastructure/db/postgres"
	"genesis/pos/infrastructure/repositories_impl"
)

// Mock de conexión de base de datos
type mockDatabaseConnectionPgx struct {
	err error
}

func (m *mockDatabaseConnectionPgx) GetDatabaseConnectionPgx() (*postgres.ConfigDatabaseConnectionPgxDriver, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Retornar nil para simular error de conexión en tests unitarios
	// Los tests de integración usarán conexiones reales
	return nil, errors.New("mock connection not fully implemented for unit tests - use integration tests")
}

// Nota: Los tests unitarios completos de repositorios con PostgreSQL requieren
// mockear pgx.Conn que es complejo. Para tests unitarios completos, se recomienda
// usar tests de integración con testcontainers o una base de datos de test.
// Este test verifica el manejo de errores de conexión.

func TestGetCurrentTimeWithPostgres_ExecuteImpl_DatabaseConnectionError(t *testing.T) {
	// Arrange
	expectedError := errors.New("connection error")
	mockDB := &mockDatabaseConnectionPgx{
		err: expectedError,
	}
	repo := &repositories_impl.GetCurrentTimeWithPostgresImpl{
		DB: mockDB,
	}

	// Act
	result, err := repo.ExecuteImpl()

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

