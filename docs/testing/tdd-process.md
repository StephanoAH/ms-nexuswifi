# Proceso TDD (Test-Driven Development)

## Introducción

Este documento describe el proceso TDD adaptado a la arquitectura hexagonal del proyecto. Aunque el proyecto actualmente no tiene tests, este proceso guiará la implementación de nuevas funcionalidades.

## Ciclo TDD: Red-Green-Refactor

El proceso TDD sigue tres pasos iterativos:

```
1. 🔴 RED: Escribir un test que falle
2. 🟢 GREEN: Escribir el código mínimo para que pase
3. 🔵 REFACTOR: Mejorar el código manteniendo los tests verdes
```

## Flujo Completo para una Nueva Funcionalidad

### Paso 1: Escribir Tests Primero

Antes de implementar cualquier funcionalidad, escribir los tests.

**Ejemplo: Nuevo caso de uso "GetUser"**

```go
// application/use_case/get_user_use_case_test.go
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

### Paso 2: Implementar Código Mínimo

Implementar solo lo necesario para que los tests pasen.

```go
// application/use_case/get_user_use_case.go
package use_case

import (
    "genesis/pos/domain/entities"
    "genesis/pos/domain/ports/infrastructure"
)

type GetUserUseCaseImpl struct {
    GetUserRepo infrastructure.GetUserRepository
}

func (g *GetUserUseCaseImpl) ExecuteUseCase() (entities.User, error) {
    user, err := g.GetUserRepo.ExecuteImpl()
    return user, err
}
```

### Paso 3: Refactorizar

Mejorar el código manteniendo los tests verdes.

```go
// Mejoras posibles:
// - Agregar validaciones
// - Mejorar manejo de errores
// - Optimizar consultas
// - Agregar logging si es necesario
```

## Orden de Testing por Capa

### 1. Domain (Entidades y Value Objects)

**Prioridad**: Alta (base de todo)

**Qué testear**:
- Lógica de negocio en entidades
- Validaciones de value objects
- Métodos de entidades

**Ejemplo**:
```go
// domain/entities/user_test.go
func TestUser_IsActive(t *testing.T) {
    user := entities.User{
        Status:    "active",
        DeletedAt: time.Time{},
    }
    
    if !user.IsActive() {
        t.Error("Expected user to be active")
    }
}
```

### 2. Infrastructure (Repositorios)

**Prioridad**: Alta (acceso a datos)

**Qué testear**:
- Queries correctas
- Mapeo de resultados
- Manejo de errores de DB
- Cierre de recursos

**Ejemplo**:
```go
// infrastructure/repositories_impl/get_user_with_postgres_test.go
func TestGetUserWithPostgres_ExecuteImpl_Success(t *testing.T) {
    // Usar testcontainers o mocks de DB
    // ...
}
```

### 3. Application (Casos de Uso)

**Prioridad**: Alta (lógica de aplicación)

**Qué testear**:
- Orquestación correcta
- Propagación de errores
- Transformación de datos

**Ejemplo**: Ver ejemplo completo arriba.

### 4. Application (Servicios)

**Prioridad**: Media-Alta (orquestación)

**Qué testear**:
- Coordinación de casos de uso
- Construcción de respuestas
- Manejo de errores

**Ejemplo**:
```go
// application/service_application/get_user_service_test.go
func TestGetUserService_ExecuteService_Success(t *testing.T) {
    // Mock del caso de uso
    // Test de orquestación
    // ...
}
```

### 5. Presentation (Handlers)

**Prioridad**: Media (integración HTTP)

**Qué testear**:
- Binding de requests
- Llamadas a servicios
- Construcción de respuestas HTTP
- Códigos de estado

**Ejemplo**:
```go
// presentation/handlers/get_user_test.go
func TestGetUser_Success(t *testing.T) {
    // Mock del servicio
    // Test de handler HTTP
    // ...
}
```

## Estrategia de Mocking

### Mocks Manuales

**Patrón**: Crear structs que implementan interfaces.

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

**Opción 1**: En el mismo archivo de test (para mocks simples)

```go
// application/use_case/get_user_use_case_test.go
package use_case

type mockGetUserRepository struct {
    // ...
}
```

**Opción 2**: En archivo separado `_mock.go` (para mocks complejos)

```go
// application/use_case/get_user_use_case_mock.go
package use_case

type mockGetUserRepository struct {
    // ...
}
```

## Estructura de Tests

### Naming de Archivos de Test

**Patrón**: `[nombre_archivo]_test.go`

```
get_user_use_case.go → get_user_use_case_test.go
get_user_repository.go → get_user_repository_test.go
```

### Naming de Funciones de Test

**Patrón**: `Test[Struct]_[Method]_[Scenario]`

```go
TestGetUserUseCase_ExecuteUseCase_Success
TestGetUserUseCase_ExecuteUseCase_Error
TestGetUserUseCase_ExecuteUseCase_NotFound
```

### Estructura AAA (Arrange-Act-Assert)

```go
func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Arrange: Preparar datos y mocks
    expectedUser := entities.User{ID: 1, Username: "test"}
    mockRepo := &mockGetUserRepository{user: expectedUser}
    useCase := &GetUserUseCaseImpl{GetUserRepo: mockRepo}
    
    // Act: Ejecutar la función a testear
    result, err := useCase.ExecuteUseCase()
    
    // Assert: Verificar resultados
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if result.ID != expectedUser.ID {
        t.Errorf("Expected ID %d, got %d", expectedUser.ID, result.ID)
    }
}
```

## Tests de Integración

### Cuándo Usar

- Tests que requieren base de datos real
- Tests de flujo completo
- Tests de performance

### Herramientas

- **testcontainers**: Para bases de datos en contenedores
- **Docker**: Para servicios externos
- **Postgres/Mongo**: Instancias de test

### Ejemplo

```go
// infrastructure/repositories_impl/get_user_with_postgres_integration_test.go
// +build integration

package repositories_impl

func TestGetUserWithPostgres_ExecuteImpl_Integration(t *testing.T) {
    // Setup testcontainers
    // Test con DB real
    // ...
}
```

## Cobertura de Tests

### Objetivo

- **Mínimo**: 70% de cobertura
- **Ideal**: 80%+ de cobertura
- **Crítico**: 90%+ para lógica de negocio

### Verificar Cobertura

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Flujo de Trabajo TDD

### 1. Requerimiento

Recibir un requerimiento claro.

**Ejemplo**: "Necesito un endpoint para obtener un usuario por ID"

### 2. Escribir Tests

Escribir tests que definan el comportamiento esperado.

```go
func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Test del caso de uso
}

func TestGetUserHandler_Success(t *testing.T) {
    // Test del handler
}
```

### 3. Implementar Mínimo

Implementar solo lo necesario para pasar los tests.

### 4. Refactorizar

Mejorar el código manteniendo tests verdes.

### 5. Iterar

Repetir para casos edge y mejoras.

## Reglas TDD

### ✅ DEBE

1. **Escribir tests primero** antes de implementar
2. **Un test por escenario** (éxito, error, edge cases)
3. **Tests independientes** (no dependen de orden)
4. **Tests rápidos** (unit tests < 1 segundo)
5. **Tests descriptivos** (nombres claros)

### ❌ NO DEBE

1. **Escribir código sin tests**
2. **Tests que dependen de otros tests**
3. **Tests lentos** (evitar I/O real en unit tests)
4. **Tests complejos** (un test, un escenario)
5. **Saltarse el refactor** (código debe mejorar)

## Ejemplo Completo: Nueva Funcionalidad

### Requerimiento

Crear endpoint para obtener usuario por ID.

### Paso 1: Tests del Caso de Uso

```go
// application/use_case/get_user_use_case_test.go
// (Ver ejemplo completo arriba)
```

### Paso 2: Implementar Caso de Uso

```go
// application/use_case/get_user_use_case.go
// (Ver ejemplo completo arriba)
```

### Paso 3: Tests del Repositorio

```go
// infrastructure/repositories_impl/get_user_with_postgres_test.go
// ...
```

### Paso 4: Implementar Repositorio

```go
// infrastructure/repositories_impl/get_user_with_postgres.go
// ...
```

### Paso 5: Tests del Handler

```go
// presentation/handlers/get_user_test.go
// ...
```

### Paso 6: Implementar Handler

```go
// presentation/handlers/get_user.go
// ...
```

### Paso 7: Refactorizar

Mejorar código manteniendo tests verdes.

## Resumen

1. **Red**: Escribir test que falle
2. **Green**: Implementar mínimo para pasar
3. **Refactor**: Mejorar código
4. **Iterar**: Repetir para nuevos casos

**Orden de testing**: Domain → Infrastructure → Application → Presentation

**Mocks**: Manuales, en archivos de test o separados

**Cobertura**: Mínimo 70%, ideal 80%+

