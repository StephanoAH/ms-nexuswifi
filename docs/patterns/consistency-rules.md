# Reglas de Consistencia

## Introducción

Este documento establece reglas para mantener la consistencia del código en todo el proyecto. Estas reglas deben seguirse estrictamente en nuevas implementaciones.

## Reglas de Arquitectura

### 1. Dependencias Entre Capas

**Regla**: Las dependencias siempre apuntan hacia adentro (hacia el dominio).

**✅ Correcto**:
```
Presentation → Application → Domain ← Infrastructure
```

**❌ Incorrecto**:
```
Domain → Infrastructure  // ❌ Domain no debe depender de Infrastructure
Application → Presentation  // ❌ Application no debe depender de Presentation
```

**Verificación**: Revisar imports en cada archivo.

### 2. Puertos en Domain

**Regla**: Todas las interfaces (puertos) deben estar en `domain/ports/`.

**✅ Correcto**:
- `domain/ports/application/use_case/get_user_use_case.go`
- `domain/ports/infrastructure/get_user_repository.go`

**❌ Incorrecto**:
- `application/use_case/get_user_use_case_interface.go`  // ❌ Interfaces fuera de domain
- `infrastructure/repositories_impl/get_user_repository_interface.go`  // ❌ Interfaces fuera de domain

### 3. Implementaciones Fuera de Domain

**Regla**: Las implementaciones NO deben estar en `domain/`.

**✅ Correcto**:
- `application/use_case/get_user_use_case.go` (implementación)
- `infrastructure/repositories_impl/get_user_with_postgres.go` (implementación)

**❌ Incorrecto**:
- `domain/use_case/get_user_use_case_impl.go`  // ❌ Implementación en domain

## Reglas de Naming

### 4. Interfaces

**Regla**: Interfaces sin sufijo (excepto cuando hay ambigüedad).

**✅ Correcto**:
```go
type GetUserUseCase interface { ... }
type GetUserRepository interface { ... }
```

**❌ Incorrecto**:
```go
type GetUserUseCaseInterface interface { ... }  // ❌ Sufijo innecesario
type GetUserRepositoryInterface interface { ... }  // ❌ Sufijo innecesario
```

**Excepción**: Cuando hay ambigüedad, usar sufijo descriptivo:
```go
type BasicAuthRepositoryPort interface { ... }  // ✅ OK: claridad
type DatabaseConnectionPgxInterface interface { ... }  // ✅ OK: claridad
```

### 5. Implementaciones

**Regla**: Implementaciones con sufijo `Impl`.

**✅ Correcto**:
```go
type GetUserUseCaseImpl struct { ... }
type GetUserServiceImpl struct { ... }
type GetUserWithPostgresImpl struct { ... }
```

**❌ Incorrecto**:
```go
type GetUserUseCase struct { ... }  // ❌ Sin sufijo
type GetUserUseCaseImplementation struct { ... }  // ❌ Sufijo incorrecto
```

### 6. Repositorios con Base de Datos

**Regla**: Indicar la tecnología en el nombre.

**✅ Correcto**:
```go
type GetUserWithPostgresImpl struct { ... }
type GetUserWithMongoImpl struct { ... }
```

**❌ Incorrecto**:
```go
type GetUserImpl struct { ... }  // ❌ No indica tecnología
type GetUserPostgresImpl struct { ... }  // ❌ Orden incorrecto
```

### 7. Métodos

**Regla**: Métodos estándar según el tipo.

**✅ Correcto**:
```go
// Casos de uso
func (uc *UseCaseImpl) ExecuteUseCase() (T, error) { ... }

// Repositorios
func (r *RepositoryImpl) ExecuteImpl() (T, error) { ... }

// Servicios
func (s *ServiceImpl) ExecuteService() T { ... }
```

**❌ Incorrecto**:
```go
func (uc *UseCaseImpl) Execute() (T, error) { ... }  // ❌ Nombre genérico
func (r *RepositoryImpl) Get() (T, error) { ... }  // ❌ No sigue patrón
```

## Reglas de Estructura

### 8. Un Concepto por Archivo

**Regla**: Un concepto principal por archivo.

**✅ Correcto**:
```
get_user_use_case.go  → GetUserUseCaseImpl
get_user_repository.go → GetUserRepository (interfaz)
get_user_with_postgres.go → GetUserWithPostgresImpl
```

**❌ Incorrecto**:
```
use_cases.go  → Múltiples casos de uso  // ❌ Múltiples conceptos
repositories.go → Múltiples repositorios  // ❌ Múltiples conceptos
```

### 9. Organización de Paquetes

**Regla**: Seguir la estructura de paquetes establecida.

**✅ Correcto**:
- Entidades → `domain/entities/`
- Puertos → `domain/ports/`
- Casos de uso → `application/use_case/`
- Repositorios → `infrastructure/repositories_impl/`

**❌ Incorrecto**:
- Entidades en `application/entities/`  // ❌ Ubicación incorrecta
- Casos de uso en `domain/use_case/`  // ❌ Implementación en domain

## Reglas de Código

### 10. Manejo de Errores

**Regla**: Siempre retornar `error` como último valor.

**✅ Correcto**:
```go
func ExecuteUseCase() (time.Time, error) { ... }
func ExecuteImpl() (entities.User, error) { ... }
```

**❌ Incorrecto**:
```go
func ExecuteUseCase() time.Time { ... }  // ❌ Sin error
func ExecuteImpl() (entities.User, bool) { ... }  // ❌ Bool en lugar de error
```

### 11. Logging de Errores

**Regla**: Usar `log.Println()` para logging de errores.

**✅ Correcto**:
```go
if err != nil {
    log.Println("Error ::: ", err)
    return nil, err
}
```

**❌ Incorrecto**:
```go
if err != nil {
    fmt.Println(err)  // ❌ No usar fmt para logging
    return nil, err
}
```

### 12. Cierre de Recursos

**Regla**: Usar `defer` para cerrar recursos.

**✅ Correcto**:
```go
conn := r.DB.GetDatabaseConnectionPgx().UsePgxDatabaseConnectionDriver()
defer conn.Close(context.Background())
```

**❌ Incorrecto**:
```go
conn := r.DB.GetDatabaseConnectionPgx().UsePgxDatabaseConnectionDriver()
// ... código ...
conn.Close(context.Background())  // ❌ Puede no ejecutarse si hay error
```

### 13. Tags de Serialización

**Regla**: Incluir tags `json:` en estructuras serializables.

**✅ Correcto**:
```go
type ResponseGetUser struct {
    StatusCode int    `json:"StatusCode"`
    Message    string `json:"Message"`
    Result     User   `json:"Result"`
}
```

**❌ Incorrecto**:
```go
type ResponseGetUser struct {
    StatusCode int    // ❌ Sin tag JSON
    Message    string
    Result     User
}
```

## Reglas de Validación

### 14. Schemas con Validación

**Regla**: Usar schemas con `ozzo-validation` para validación compleja.

**✅ Correcto**:
```go
type CreateUserBody struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

func (c CreateUserBody) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.Username, validation.Required),
        validation.Field(&c.Email, validation.Required, validation.Email),
    )
}
```

**❌ Incorrecto**:
```go
// Validación manual en handler
if body.Username == "" {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
    return
}
```

### 15. Middlewares de Validación

**Regla**: Usar middlewares genéricos para validación.

**✅ Correcto**:
```go
router.POST("/",
    middleware.ValidateBody(new(schema.CreateUserBody)),
    handlers.CreateUser,
)
```

**❌ Incorrecto**:
```go
// Validación manual en handler
func CreateUser(ctx *gin.Context) {
    var body schema.CreateUserBody
    if err := ctx.BindJSON(&body); err != nil {
        // ...
    }
    if err := body.Validate(); err != nil {
        // ...
    }
}
```

## Reglas de Container

### 16. Resolución de Dependencias

**Regla**: Usar funciones `ResolveXxx()` del container.

**✅ Correcto**:
```go
var GetUserRepository infrastructure.GetUserRepository

func ResolveGetUserRepository() infrastructure.GetUserRepository {
    if GetUserRepository != nil {
        return GetUserRepository
    }
    GetUserRepository = &repositories_impl.GetUserWithPostgresImpl{
        DB: DatabaseConnectionToLecWithPgx,
    }
    return GetUserRepository
}
```

**❌ Incorrecto**:
```go
// Crear instancias directamente en handlers
func CreateUser(ctx *gin.Context) {
    repo := &repositories_impl.GetUserWithPostgresImpl{...}  // ❌ No usar container
}
```

### 17. Inicialización en StartContainer

**Regla**: Inicializar todas las dependencias en `StartContainer()`.

**✅ Correcto**:
```go
func StartContainer() {
    DatabaseConnectionToLecWithPgx = ResolveDatabaseConnectionToLecWithPgx()
    GetUserRepository = ResolveGetUserRepository()
    GetUserUseCase = ResolveGetUserUseCase()
    GetUserService = ResolveGetUserService()
}
```

**❌ Incorrecto**:
```go
// Inicialización dispersa
func SomeHandler(ctx *gin.Context) {
    if GetUserService == nil {
        GetUserService = ResolveGetUserService()  // ❌ Inicialización tardía
    }
}
```

## Reglas de Testing

### 18. Naming de Tests

**Regla**: `Test[Struct]_[Method]_[Scenario]`

**✅ Correcto**:
```go
TestGetUserUseCase_ExecuteUseCase_Success
TestGetUserUseCase_ExecuteUseCase_Error
TestGetUserUseCase_ExecuteUseCase_NotFound
```

**❌ Incorrecto**:
```go
TestGetUser  // ❌ Muy genérico
TestExecute  // ❌ No indica struct ni método
test_get_user  // ❌ No sigue convención Go
```

### 19. Estructura AAA

**Regla**: Usar patrón Arrange-Act-Assert en tests.

**✅ Correcto**:
```go
func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Arrange
    mockRepo := &mockGetUserRepository{...}
    useCase := &GetUserUseCaseImpl{GetUserRepo: mockRepo}
    
    // Act
    result, err := useCase.ExecuteUseCase()
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
}
```

**❌ Incorrecto**:
```go
func TestGetUserUseCase_ExecuteUseCase_Success(t *testing.T) {
    // Mezcla de arrange, act y assert
    useCase := &GetUserUseCaseImpl{...}
    result, err := useCase.ExecuteUseCase()
    useCase.GetUserRepo = &mockGetUserRepository{...}  // ❌ Arrange después de Act
    if err != nil {
        // ...
    }
}
```

## Reglas de Documentación

### 20. Comentarios en Código

**Regla**: Comentarios cuando el código no es autoexplicativo.

**✅ Correcto**:
```go
// Check if the content type is multipart/form-data
if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
    // Log request without body if it's multipart (likely contains files)
    log.Printf("Request [multipart]: %s %s from %s", ...)
}
```

**❌ Incorrecto**:
```go
// Get user
user := getUser()  // ❌ Comentario obvio
```

## Checklist de Consistencia

Antes de hacer commit, verificar:

- [ ] Dependencias apuntan hacia adentro
- [ ] Interfaces en `domain/ports/`
- [ ] Implementaciones con sufijo `Impl`
- [ ] Métodos siguen convenciones (`ExecuteUseCase()`, `ExecuteImpl()`, `ExecuteService()`)
- [ ] Errores retornados como último valor
- [ ] Logging con `log.Println()`
- [ ] Recursos cerrados con `defer`
- [ ] Tags JSON en estructuras serializables
- [ ] Validación en schemas con `ozzo-validation`
- [ ] Dependencias resueltas en container
- [ ] Tests siguen naming `Test[Struct]_[Method]_[Scenario]`
- [ ] Tests usan patrón AAA

## Resumen

Estas reglas aseguran que el código sea:
- **Consistente**: Mismo patrón en todo el proyecto
- **Mantenible**: Fácil de entender y modificar
- **Escalable**: Fácil agregar nuevas funcionalidades
- **Testeable**: Estructura clara para tests

**Regla de Oro**: Si algo no sigue estas reglas, es una inconsistencia que debe corregirse.

