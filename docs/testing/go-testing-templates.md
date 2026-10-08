# Plantillas de Testing en Go

## Introducción

Este documento proporciona plantillas listas para usar en los tests del proyecto, siguiendo las convenciones establecidas.

## Plantilla: Test de Caso de Uso

```go
package use_case

import (
    "errors"
    "testing"
    "genesis/pos/domain/entities"
    "genesis/pos/domain/ports/infrastructure"
)

// Mock del repositorio
type mockGetUserRepository struct {
    user entities.User
    err  error
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    return m.user, m.err
}

func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Arrange
    expectedUser := entities.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }
    mockRepo := &mockGetUserRepository{
        user: expectedUser,
        err:  nil,
    }
    useCase := &GetUserUseCaseImpl{
        GetUserRepo: mockRepo,
    }
    
    // Act
    result, err := useCase.ExecuteUseCase()
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if result.ID != expectedUser.ID {
        t.Errorf("Expected ID %d, got %d", expectedUser.ID, result.ID)
    }
    if result.Username != expectedUser.Username {
        t.Errorf("Expected Username %s, got %s", expectedUser.Username, result.Username)
    }
}

func TestGetUserUseCase_ExecuteUseCase_Error(t *testing.T) {
    // Arrange
    expectedError := errors.New("database error")
    mockRepo := &mockGetUserRepository{
        user: entities.User{},
        err:  expectedError,
    }
    useCase := &GetUserUseCaseImpl{
        GetUserRepo: mockRepo,
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
    if result.ID != 0 {
        t.Errorf("Expected zero value, got %d", result.ID)
    }
}
```

## Plantilla: Test de Repositorio (Unit Test con Mock de DB)

```go
package repositories_impl

import (
    "context"
    "errors"
    "testing"
    "genesis/pos/domain/entities"
    "genesis/pos/domain/ports/infrastructure"
)

// Mock de conexión de base de datos
type mockDatabaseConnectionPgx struct {
    conn *mockPgxConn
}

func (m *mockDatabaseConnectionPgx) GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver {
    return &postgres.ConfigDatabaseConnectionPgxDriver{
        PgxConn: m.conn,
    }
}

type mockPgxConn struct {
    queryResult entities.User
    queryErr    error
}

func (m *mockPgxConn) QueryRow(ctx context.Context, sql string, args ...interface{}) *mockRow {
    return &mockRow{
        user: m.queryResult,
        err:  m.queryErr,
    }
}

type mockRow struct {
    user entities.User
    err  error
}

func (m *mockRow) Scan(dest ...interface{}) error {
    if m.err != nil {
        return m.err
    }
    // Simular scan (simplificado)
    return nil
}

func TestGetUserWithPostgres_ExecuteImpl_Success(t *testing.T) {
    // Arrange
    expectedUser := entities.User{
        ID:       1,
        Username: "testuser",
    }
    mockDB := &mockDatabaseConnectionPgx{
        conn: &mockPgxConn{
            queryResult: expectedUser,
            queryErr:    nil,
        },
    }
    repo := &GetUserWithPostgresImpl{
        DB: mockDB,
    }
    
    // Act
    result, err := repo.ExecuteImpl()
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if result.ID != expectedUser.ID {
        t.Errorf("Expected ID %d, got %d", expectedUser.ID, result.ID)
    }
}
```

## Plantilla: Test de Repositorio (Integración)

```go
// +build integration

package repositories_impl

import (
    "context"
    "testing"
    "genesis/pos/domain/entities"
    "genesis/pos/infrastructure/db/postgres/drivers"
)

func TestGetUserWithPostgres_ExecuteImpl_Integration(t *testing.T) {
    // Arrange: Setup testcontainers o DB de test
    ctx := context.Background()
    
    // Configurar conexión de test
    configConnection := &drivers.ConfigConnectionPgx{
        UrlToConnect: getTestDatabaseURL(), // Función helper
    }
    
    dbConnection := configConnection.GetDatabaseConnectionPgx()
    repo := &GetUserWithPostgresImpl{
        DB: dbConnection,
    }
    
    // Setup: Insertar datos de test
    setupTestData(t, ctx, dbConnection)
    
    // Act
    result, err := repo.ExecuteImpl()
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if result.ID == 0 {
        t.Error("Expected non-zero ID")
    }
    
    // Cleanup
    cleanupTestData(t, ctx, dbConnection)
}

func getTestDatabaseURL() string {
    // Retornar URL de DB de test
    return os.Getenv("TEST_DATABASE_URL")
}

func setupTestData(t *testing.T, ctx context.Context, db *postgres.ConfigDatabaseConnectionPgxDriver) {
    // Insertar datos de test
}

func cleanupTestData(t *testing.T, ctx context.Context, db *postgres.ConfigDatabaseConnectionPgxDriver) {
    // Limpiar datos de test
}
```

## Plantilla: Test de Servicio

```go
package service_application

import (
    "errors"
    "testing"
    "genesis/pos/domain/entities"
    "genesis/pos/domain/ports/application/use_case"
)

// Mock del caso de uso
type mockGetUserUseCase struct {
    user entities.User
    err  error
}

func (m *mockGetUserUseCase) ExecuteUseCase() (entities.User, error) {
    return m.user, m.err
}

func TestGetUserService_ExecuteService_Success(t *testing.T) {
    // Arrange
    expectedUser := entities.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }
    mockUseCase := &mockGetUserUseCase{
        user: expectedUser,
        err:  nil,
    }
    service := &GetUserServiceImpl{
        GetUserUseCase: mockUseCase,
    }
    
    // Act
    response := service.ExecuteService()
    
    // Assert
    if response.StatusCode != 200 {
        t.Errorf("Expected StatusCode 200, got %d", response.StatusCode)
    }
    if response.Result.ID != expectedUser.ID {
        t.Errorf("Expected User ID %d, got %d", expectedUser.ID, response.Result.ID)
    }
}

func TestGetUserService_ExecuteService_Error(t *testing.T) {
    // Arrange
    expectedError := errors.New("database error")
    mockUseCase := &mockGetUserUseCase{
        user: entities.User{},
        err:  expectedError,
    }
    service := &GetUserServiceImpl{
        GetUserUseCase: mockUseCase,
    }
    
    // Act
    response := service.ExecuteService()
    
    // Assert
    if response.StatusCode != 500 {
        t.Errorf("Expected StatusCode 500, got %d", response.StatusCode)
    }
    if response.Message == "" {
        t.Error("Expected error message")
    }
}
```

## Plantilla: Test de Handler

```go
package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "genesis/pos/domain/entities"
    "genesis/pos/presentation/container"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// Mock del servicio
type mockGetUserService struct {
    response entities.ResponseGetUser
}

func (m *mockGetUserService) ExecuteService() entities.ResponseGetUser {
    return m.response
}

func TestGetUser_Success(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    
    expectedResponse := entities.ResponseGetUser{
        StatusCode: 200,
        Message:    "User retrieved successfully",
        Result: entities.User{
            ID:       1,
            Username: "testuser",
        },
    }
    
    // Mock del servicio en el container
    originalService := container.GetUserService
    container.GetUserService = &mockGetUserService{
        response: expectedResponse,
    }
    defer func() {
        container.GetUserService = originalService
    }()
    
    // Crear request
    req, _ := http.NewRequest("GET", "/api/users/1", nil)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
    
    // Act
    GetUser(c)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response entities.ResponseHttp
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response.Success)
    assert.Equal(t, 200, response.Code)
}

func TestGetUser_Error(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    
    errorResponse := entities.ResponseGetUser{
        StatusCode: 500,
        Message:    "Error retrieving user",
        Result:     entities.User{},
    }
    
    originalService := container.GetUserService
    container.GetUserService = &mockGetUserService{
        response: errorResponse,
    }
    defer func() {
        container.GetUserService = originalService
    }()
    
    req, _ := http.NewRequest("GET", "/api/users/1", nil)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
    
    // Act
    GetUser(c)
    
    // Assert
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var response entities.ResponseHttp
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.False(t, response.Success)
}
```

## Plantilla: Table-Driven Test

```go
package use_case

import (
    "errors"
    "testing"
    "genesis/pos/domain/entities"
)

func TestGetUserUseCase_ExecuteUseCase(t *testing.T) {
    tests := []struct {
        name     string
        mockUser entities.User
        mockErr  error
        wantUser entities.User
        wantErr  bool
    }{
        {
            name: "success",
            mockUser: entities.User{
                ID:       1,
                Username: "testuser",
                Email:    "test@example.com",
            },
            mockErr: nil,
            wantUser: entities.User{
                ID:       1,
                Username: "testuser",
                Email:    "test@example.com",
            },
            wantErr: false,
        },
        {
            name:     "database error",
            mockUser: entities.User{},
            mockErr:  errors.New("database error"),
            wantUser: entities.User{},
            wantErr:  true,
        },
        {
            name:     "not found",
            mockUser: entities.User{},
            mockErr:  errors.New("user not found"),
            wantUser: entities.User{},
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := &mockGetUserRepository{
                user: tt.mockUser,
                err:  tt.mockErr,
            }
            useCase := &GetUserUseCaseImpl{
                GetUserRepo: mockRepo,
            }
            
            // Act
            result, err := useCase.ExecuteUseCase()
            
            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("ExecuteUseCase() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr {
                if result.ID != tt.wantUser.ID {
                    t.Errorf("ExecuteUseCase() ID = %d, want %d", result.ID, tt.wantUser.ID)
                }
                if result.Username != tt.wantUser.Username {
                    t.Errorf("ExecuteUseCase() Username = %s, want %s", result.Username, tt.wantUser.Username)
                }
            }
        })
    }
}
```

## Plantilla: Test Helper

```go
package use_case

import "genesis/pos/domain/entities"

// Helper para crear mock de repositorio
func createMockUserRepository(user entities.User, err error) *mockGetUserRepository {
    return &mockGetUserRepository{
        user: user,
        err:  err,
    }
}

// Helper para crear usuario de test
func createTestUser() entities.User {
    return entities.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }
}

// Helper para crear caso de uso de test
func createTestUseCase(repo *mockGetUserRepository) *GetUserUseCaseImpl {
    return &GetUserUseCaseImpl{
        GetUserRepo: repo,
    }
}
```

## Resumen de Plantillas

1. **Test de Caso de Uso**: Mock de repositorio, test de éxito y error
2. **Test de Repositorio**: Mock de DB o test de integración
3. **Test de Servicio**: Mock de caso de uso, test de orquestación
4. **Test de Handler**: Mock de servicio, test de HTTP
5. **Table-Driven Test**: Múltiples escenarios en un test
6. **Test Helpers**: Funciones auxiliares para tests

Cada plantilla sigue el patrón AAA (Arrange-Act-Assert) y las convenciones de naming del proyecto.

