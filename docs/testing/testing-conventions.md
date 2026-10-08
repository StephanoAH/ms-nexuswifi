# Convenciones de Testing

## Introducción

Este documento describe las convenciones de testing específicas del proyecto, basadas en las mejores prácticas de Go y adaptadas a la arquitectura hexagonal.

## Herramientas de Testing

### Testing Estándar

El proyecto usa el paquete `testing` estándar de Go.

```go
import "testing"
```

### Librerías Adicionales (Recomendadas)

Aunque no están en el proyecto actual, se recomiendan:

- **testify**: Para assertions y mocks
- **testcontainers**: Para tests de integración con DB
- **gomock**: Para generación de mocks (opcional)

## Estructura de Archivos de Test

### Naming

**Patrón**: `[nombre_archivo]_test.go`

```
get_user_use_case.go → get_user_use_case_test.go
get_user_repository.go → get_user_repository_test.go
get_user_handler.go → get_user_handler_test.go
```

### Ubicación

Los archivos de test están en el **mismo paquete** que el código a testear.

```
application/use_case/
├── get_user_use_case.go
└── get_user_use_case_test.go  ✅ Mismo paquete
```

**Excepción**: Tests de integración pueden estar en subdirectorio `integration/` o usar build tags.

## Naming de Funciones de Test

### Patrón

**Formato**: `Test[Struct]_[Method]_[Scenario]`

### Ejemplos

```go
// Caso exitoso
TestGetUserUseCase_ExecuteUseCase_Success

// Caso de error
TestGetUserUseCase_ExecuteUseCase_Error

// Caso específico
TestGetUserUseCase_ExecuteUseCase_NotFound

// Caso edge
TestGetUserUseCase_ExecuteUseCase_EmptyResult
```

### Reglas

1. **PascalCase** para nombres de test
2. **Descriptivo**: Indica qué se testea y el escenario
3. **Específico**: Un escenario por test

## Estructura de Tests: AAA Pattern

### Arrange-Act-Assert

```go
func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Arrange: Preparar datos, mocks, configuración
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
    
    // Act: Ejecutar la función a testear
    result, err := useCase.ExecuteUseCase()
    
    // Assert: Verificar resultados
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
```

### Comentarios en Tests

**Opcional pero recomendado**: Comentarios para clarificar secciones AAA.

```go
func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Arrange
    // ...
    
    // Act
    // ...
    
    // Assert
    // ...
}
```

## Mocks

### Mocks Manuales

**Patrón**: Struct que implementa la interfaz.

```go
type mockGetUserRepository struct {
    user entities.User
    err  error
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    return m.user, m.err
}
```

### Ubicación de Mocks

**Opción 1**: En el mismo archivo de test (recomendado para mocks simples)

```go
// application/use_case/get_user_use_case_test.go
package use_case

type mockGetUserRepository struct {
    user entities.User
    err  error
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    return m.user, m.err
}

func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Usar mock
}
```

**Opción 2**: En archivo separado `_mock.go` (para mocks complejos o reutilizables)

```go
// application/use_case/get_user_use_case_mock.go
package use_case

type mockGetUserRepository struct {
    // Implementación compleja
}
```

### Mocks con Comportamiento

Para mocks más sofisticados, agregar métodos de configuración:

```go
type mockGetUserRepository struct {
    user entities.User
    err  error
    callCount int
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    m.callCount++
    return m.user, m.err
}

func (m *mockGetUserRepository) WithUser(user entities.User) *mockGetUserRepository {
    m.user = user
    return m
}

func (m *mockGetUserRepository) WithError(err error) *mockGetUserRepository {
    m.err = err
    return m
}
```

## Table-Driven Tests

### Patrón

Para múltiples escenarios similares, usar table-driven tests.

```go
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
            mockUser: entities.User{ID: 1, Username: "test"},
            mockErr:  nil,
            wantUser: entities.User{ID: 1, Username: "test"},
            wantErr:  false,
        },
        {
            name: "error",
            mockUser: entities.User{},
            mockErr:  errors.New("database error"),
            wantUser: entities.User{},
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockGetUserRepository{
                user: tt.mockUser,
                err:  tt.mockErr,
            }
            useCase := &GetUserUseCaseImpl{
                GetUserRepo: mockRepo,
            }
            
            result, err := useCase.ExecuteUseCase()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ExecuteUseCase() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result.ID != tt.wantUser.ID {
                t.Errorf("ExecuteUseCase() ID = %d, want %d", result.ID, tt.wantUser.ID)
            }
        })
    }
}
```

## Assertions

### Testing Estándar

**Patrón**: Usar `if` statements con `t.Errorf` o `t.Fatalf`.

```go
if err != nil {
    t.Fatalf("Expected no error, got %v", err)
}

if result.ID != expectedID {
    t.Errorf("Expected ID %d, got %d", expectedID, result.ID)
}
```

### Con testify (Recomendado)

Si se agrega `testify`, usar assertions:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Arrange
    // ...
    
    // Act
    result, err := useCase.ExecuteUseCase()
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, expectedUser.ID, result.ID)
    assert.Equal(t, expectedUser.Username, result.Username)
}
```

**Diferencia**:
- `require.*`: Falla el test inmediatamente (para precondiciones)
- `assert.*`: Continúa el test (para verificaciones)

## Tests de Integración

### Build Tags

Usar build tags para separar tests de integración.

```go
// +build integration

package repositories_impl

func TestGetUserWithPostgres_ExecuteImpl_Integration(t *testing.T) {
    // Test con DB real
}
```

### Ejecutar Tests de Integración

```bash
go test -tags=integration ./...
```

### Testcontainers

Para tests con bases de datos reales:

```go
// infrastructure/repositories_impl/get_user_with_postgres_integration_test.go
// +build integration

package repositories_impl

import (
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetUserWithPostgres_ExecuteImpl_Integration(t *testing.T) {
    // Setup testcontainers
    ctx := context.Background()
    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil {
        t.Fatalf("Failed to start container: %v", err)
    }
    defer postgresContainer.Terminate(ctx)
    
    // Test con DB real
    // ...
}
```

## Helpers y Utilidades

### Test Helpers

Para código reutilizable en tests:

```go
// application/use_case/test_helpers.go
package use_case

func createMockUserRepository(user entities.User, err error) *mockGetUserRepository {
    return &mockGetUserRepository{
        user: user,
        err:  err,
    }
}

func createTestUser() entities.User {
    return entities.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }
}
```

## Cobertura

### Verificar Cobertura

```bash
# Cobertura de un paquete
go test -cover ./application/use_case

# Cobertura completa
go test -cover ./...

# Cobertura con perfil
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Objetivos de Cobertura

- **Mínimo**: 70%
- **Ideal**: 80%+
- **Crítico**: 90%+ para lógica de negocio

## Ejecutar Tests

### Todos los Tests

```bash
go test ./...
```

### Tests de un Paquete

```bash
go test ./application/use_case
```

### Tests Específicos

```bash
go test -run TestGetUserUseCase_ExecuteUseCase_Success ./application/use_case
```

### Tests con Verbose

```bash
go test -v ./...
```

### Tests en Paralelo

```bash
go test -parallel 4 ./...
```

## Reglas de Testing

### ✅ DEBE

1. **Un test por escenario** (éxito, error, edge cases)
2. **Tests independientes** (no dependen de orden)
3. **Tests rápidos** (unit tests < 1 segundo)
4. **Tests descriptivos** (nombres claros)
5. **Mocks simples** (solo lo necesario)
6. **AAA pattern** (Arrange-Act-Assert)

### ❌ NO DEBE

1. **Tests que dependen de otros tests**
2. **Tests lentos** (evitar I/O real en unit tests)
3. **Tests complejos** (un test, un escenario)
4. **Mocks excesivamente complejos**
5. **Tests sin assertions claras**

## Resumen

- **Naming**: `Test[Struct]_[Method]_[Scenario]`
- **Estructura**: AAA (Arrange-Act-Assert)
- **Mocks**: Manuales, en mismo archivo o separados
- **Table-driven**: Para múltiples escenarios
- **Integración**: Con build tags y testcontainers
- **Cobertura**: Mínimo 70%, ideal 80%+

