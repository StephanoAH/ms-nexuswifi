# Estrategia de Mocking

## Introducción

Este documento describe la estrategia de mocking para tests en el proyecto, basada en la arquitectura hexagonal y las convenciones detectadas.

## Principio Fundamental

**Regla**: Mockear dependencias externas, no el código que se está testeando.

En arquitectura hexagonal:
- Mockear **puertos** (interfaces), no implementaciones
- Mockear **repositorios** en tests de casos de uso
- Mockear **casos de uso** en tests de servicios
- Mockear **servicios** en tests de handlers

## Tipos de Mocks

### 1. Mocks Manuales

**Uso**: Para la mayoría de los tests del proyecto.

**Ventajas**:
- Simple y directo
- Sin dependencias externas
- Control total sobre el comportamiento

**Desventajas**:
- Código repetitivo
- Mantenimiento manual

**Ejemplo**:
```go
type mockGetUserRepository struct {
    user entities.User
    err  error
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    return m.user, m.err
}
```

### 2. Mocks con Comportamiento

**Uso**: Cuando se necesita verificar llamadas o contar invocaciones.

**Ejemplo**:
```go
type mockGetUserRepository struct {
    user      entities.User
    err       error
    callCount int
    lastCall  time.Time
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    m.callCount++
    m.lastCall = time.Now()
    return m.user, m.err
}

func (m *mockGetUserRepository) GetCallCount() int {
    return m.callCount
}
```

### 3. Mocks con Configuración Fluent

**Uso**: Para tests más legibles y configurables.

**Ejemplo**:
```go
type mockGetUserRepository struct {
    user entities.User
    err  error
}

func NewMockGetUserRepository() *mockGetUserRepository {
    return &mockGetUserRepository{}
}

func (m *mockGetUserRepository) WithUser(user entities.User) *mockGetUserRepository {
    m.user = user
    return m
}

func (m *mockGetUserRepository) WithError(err error) *mockGetUserRepository {
    m.err = err
    return m
}

// Uso en test
mockRepo := NewMockGetUserRepository().
    WithUser(entities.User{ID: 1, Username: "test"}).
    WithError(nil)
```

## Estrategia por Capa

### Domain Layer

**No se mockea**: Las entidades y value objects son datos puros.

**Tests directos**:
```go
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

### Application Layer - Casos de Uso

**Mockear**: Repositorios (puertos de infraestructura)

**Ejemplo**:
```go
// Mock del repositorio
type mockGetUserRepository struct {
    user entities.User
    err  error
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    return m.user, m.err
}

func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    mockRepo := &mockGetUserRepository{
        user: entities.User{ID: 1, Username: "test"},
        err:  nil,
    }
    
    useCase := &GetUserUseCaseImpl{
        GetUserRepo: mockRepo,  // ✅ Mock del puerto
    }
    
    result, err := useCase.ExecuteUseCase()
    // Assertions...
}
```

### Application Layer - Servicios

**Mockear**: Casos de uso (puertos de aplicación)

**Ejemplo**:
```go
// Mock del caso de uso
type mockGetUserUseCase struct {
    user entities.User
    err  error
}

func (m *mockGetUserUseCase) ExecuteUseCase() (entities.User, error) {
    return m.user, m.err
}

func TestGetUserService_ExecuteService_Success(t *testing.T) {
    mockUseCase := &mockGetUserUseCase{
        user: entities.User{ID: 1, Username: "test"},
        err:  nil,
    }
    
    service := &GetUserServiceImpl{
        GetUserUseCase: mockUseCase,  // ✅ Mock del caso de uso
    }
    
    response := service.ExecuteService()
    // Assertions...
}
```

### Infrastructure Layer - Repositorios

**Opciones**:

1. **Mock de conexión de DB** (para unit tests):
```go
type mockDatabaseConnectionPgx struct {
    conn *mockPgxConn
}

func (m *mockDatabaseConnectionPgx) GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver {
    return &postgres.ConfigDatabaseConnectionPgxDriver{
        PgxConn: m.conn,
    }
}
```

2. **Test de integración** (para tests con DB real):
```go
// +build integration

func TestGetUserWithPostgres_ExecuteImpl_Integration(t *testing.T) {
    // Usar testcontainers o DB de test
}
```

### Presentation Layer - Handlers

**Mockear**: Servicios (del container)

**Ejemplo**:
```go
func TestGetUser_Success(t *testing.T) {
    // Guardar servicio original
    originalService := container.GetUserService
    
    // Mock del servicio
    mockService := &mockGetUserService{
        response: entities.ResponseGetUser{
            StatusCode: 200,
            Result:     entities.User{ID: 1},
        },
    }
    
    // Reemplazar en container
    container.GetUserService = mockService
    defer func() {
        container.GetUserService = originalService  // Restaurar
    }()
    
    // Test del handler
    // ...
}
```

## Ubicación de Mocks

### Opción 1: En el mismo archivo de test (Recomendado)

**Para**: Mocks simples y específicos de un test.

```go
// application/use_case/get_user_use_case_test.go
package use_case

type mockGetUserRepository struct {
    // ...
}

func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Usar mock
}
```

### Opción 2: En archivo separado `_mock.go`

**Para**: Mocks complejos o reutilizables.

```go
// application/use_case/get_user_use_case_mock.go
package use_case

type mockGetUserRepository struct {
    // Implementación compleja
}

type mockGetUserRepositoryBuilder struct {
    // Builder pattern para configuración
}
```

### Opción 3: En paquete `test` o `mocks`

**Para**: Mocks compartidos entre múltiples paquetes.

```
test/
├── mocks/
│   ├── repository_mocks.go
│   └── use_case_mocks.go
```

**Nota**: Esta opción no se usa actualmente en el proyecto, pero es válida para mocks muy reutilizables.

## Verificación de Llamadas

### Contar Invocaciones

```go
type mockGetUserRepository struct {
    user      entities.User
    err       error
    callCount int
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    m.callCount++
    return m.user, m.err
}

func TestGetUserUseCase_CallsRepository(t *testing.T) {
    mockRepo := &mockGetUserRepository{}
    useCase := &GetUserUseCaseImpl{GetUserRepo: mockRepo}
    
    useCase.ExecuteUseCase()
    
    if mockRepo.callCount != 1 {
        t.Errorf("Expected 1 call, got %d", mockRepo.callCount)
    }
}
```

### Verificar Parámetros

```go
type mockGetUserRepository struct {
    user        entities.User
    err         error
    receivedID  int64
}

func (m *mockGetUserRepository) GetUserByID(id int64) (entities.User, error) {
    m.receivedID = id
    return m.user, m.err
}

func TestGetUserUseCase_PassesCorrectID(t *testing.T) {
    expectedID := int64(123)
    mockRepo := &mockGetUserRepository{}
    useCase := &GetUserUseCaseImpl{GetUserRepo: mockRepo}
    
    useCase.ExecuteUseCase(expectedID)
    
    if mockRepo.receivedID != expectedID {
        t.Errorf("Expected ID %d, got %d", expectedID, mockRepo.receivedID)
    }
}
```

## Mocks para Errores

### Simular Diferentes Errores

```go
func TestGetUserUseCase_HandlesDifferentErrors(t *testing.T) {
    tests := []struct {
        name    string
        mockErr error
        wantErr bool
    }{
        {
            name:    "database error",
            mockErr:  errors.New("connection failed"),
            wantErr:  true,
        },
        {
            name:    "not found",
            mockErr:  errors.New("user not found"),
            wantErr:  true,
        },
        {
            name:    "no error",
            mockErr:  nil,
            wantErr:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockGetUserRepository{err: tt.mockErr}
            useCase := &GetUserUseCaseImpl{GetUserRepo: mockRepo}
            
            _, err := useCase.ExecuteUseCase()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Expected error %v, got %v", tt.wantErr, err != nil)
            }
        })
    }
}
```

## Mocks para Timeouts y Context

### Simular Timeouts

```go
type mockGetUserRepository struct {
    user     entities.User
    err      error
    delay    time.Duration
}

func (m *mockGetUserRepository) ExecuteImpl() (entities.User, error) {
    if m.delay > 0 {
        time.Sleep(m.delay)
    }
    return m.user, m.err
}

func TestGetUserUseCase_HandlesTimeout(t *testing.T) {
    mockRepo := &mockGetUserRepository{
        delay: 2 * time.Second,
    }
    useCase := &GetUserUseCaseImpl{GetUserRepo: mockRepo}
    
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    // Test con timeout
    // ...
}
```

## Resumen de Estrategia

### Por Capa

| Capa | Qué Mockear | Ejemplo |
|------|------------|---------|
| **Domain** | Nada (datos puros) | Tests directos |
| **Use Case** | Repositorios (puertos) | `mockGetUserRepository` |
| **Service** | Casos de uso (puertos) | `mockGetUserUseCase` |
| **Repository** | Conexión DB o test de integración | `mockDatabaseConnection` |
| **Handler** | Servicios (container) | `mockGetUserService` |

### Reglas

1. ✅ **Mockear puertos** (interfaces), no implementaciones
2. ✅ **Mocks simples** en el mismo archivo de test
3. ✅ **Mocks complejos** en archivo separado `_mock.go`
4. ✅ **Verificar llamadas** cuando sea relevante
5. ✅ **Tests de integración** para repositorios con DB real

### Evitar

1. ❌ Mockear el código que se está testeando
2. ❌ Mocks excesivamente complejos
3. ❌ Mocks que no reflejan el comportamiento real
4. ❌ Dependencias externas en mocks (evitar I/O real)

