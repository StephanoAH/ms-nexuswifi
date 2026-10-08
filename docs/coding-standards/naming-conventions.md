# Convenciones de Nombres

## Introducción

Este documento detalla las convenciones de nombres específicas del proyecto, basadas en los patrones detectados en el código existente.

## Resumen Rápido

| Elemento | Convención | Ejemplo |
|----------|-----------|---------|
| **Interfaces** | Sin sufijo | `GetCurrentTimeUseCase` |
| **Implementaciones** | Sufijo `Impl` | `GetCurrentTimeUseCaseImpl` |
| **Repositorios con DB** | `With[DB]Impl` | `GetCurrentTimeWithPostgresImpl` |
| **Métodos Use Case** | `ExecuteUseCase()` | `ExecuteUseCase() (T, error)` |
| **Métodos Repository** | `ExecuteImpl()` o específico | `ExecuteImpl()` o `GetCredentials()` |
| **Métodos Service** | `ExecuteService()` | `ExecuteService() T` |
| **Paquetes** | minúsculas, una palabra | `entities`, `use_case` |
| **Structs** | PascalCase | `ResponseGetCurrentTime` |
| **Variables exportadas** | PascalCase | `DatabaseConnectionToMongo` |
| **Variables no exportadas** | camelCase | `appParams` |

## Interfaces

### Regla General

**Interfaces NO llevan sufijo**, a menos que haya ambigüedad.

### Patrón Estándar

```go
// ✅ Correcto
type GetCurrentTimeUseCase interface {
    ExecuteUseCase() (time.Time, error)
}

type GetCurrentTimeRepository interface {
    ExecuteImpl() (time.Time, error)
}
```

### Excepciones (Ambigüedad)

Cuando hay ambigüedad o se necesita claridad adicional:

```go
// ✅ Correcto: "Port" para claridad
type BasicAuthRepositoryPort interface {
    GetCredentials(integrationID int64, username string) (entities.BasicAuth, error)
}

// ✅ Correcto: "Interface" para conexiones
type DatabaseConnectionPgxInterface interface {
    GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver
}
```

### Naming de Interfaces por Tipo

| Tipo | Patrón | Ejemplo |
|------|--------|---------|
| Use Case | `Get[Entity]UseCase` | `GetCurrentTimeUseCase` |
| Service | `Get[Entity]Service` | `GetCurrentTimeService` |
| Repository | `Get[Entity]Repository` | `GetCurrentTimeRepository` |
| Connection | `DatabaseConnection[Type]Interface` | `DatabaseConnectionPgxInterface` |

## Implementaciones

### Regla General

**Implementaciones llevan sufijo `Impl`**.

### Patrón Estándar

```go
// ✅ Correcto
type GetCurrentTimeUseCaseImpl struct {
    GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository
}

type GetCurrentTimeServiceImpl struct {
    GetCurrentTimeServiceUseCase use_case.GetCurrentTimeUseCase
}
```

### Repositorios con Base de Datos Específica

**Patrón**: `[Nombre]With[DB]Impl`

```go
// ✅ Correcto: PostgreSQL
type GetCurrentTimeWithPostgresImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}

// ✅ Correcto: MongoDB
type GetSalesWithMongoImpl struct {
    MongoDB infrastructure.DatabaseConnectionMongoInterface
}
```

### Repositorios Genéricos

**Patrón**: `[Nombre]RepositoryImpl`

```go
// ✅ Correcto
type BasicAuthRepositoryImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}
```

## Métodos

### Casos de Uso

**Patrón**: `ExecuteUseCase()`

```go
func (g *GetCurrentTimeUseCaseImpl) ExecuteUseCase() (time.Time, error) {
    CurrentTime, err := g.GetCurrentTimeRepo.ExecuteImpl()
    return CurrentTime, err
}
```

**Variaciones detectadas**:
- `ExecuteUseCase()`: Método genérico
- `GetCredentialsUseCase(...)`: Método específico del dominio (menos común)

### Repositorios

**Patrón Principal**: `ExecuteImpl()`

```go
func (r *GetCurrentTimeWithPostgresImpl) ExecuteImpl() (time.Time, error) {
    // Implementación
}
```

**Patrón Alternativo**: Métodos específicos del dominio

```go
func (b *BasicAuthRepositoryImpl) GetCredentials(integrationID int64, username string) (entities.BasicAuth, error) {
    // Implementación específica
}
```

**Cuándo usar cada uno**:
- `ExecuteImpl()`: Cuando el repositorio tiene una operación principal
- Método específico: Cuando el repositorio tiene múltiples operaciones distintas

### Servicios

**Patrón**: `ExecuteService()`

```go
func (g *GetCurrentTimeServiceImpl) ExecuteService() entities.ResponseGetCurrentTime {
    // Orquestación de casos de uso
}
```

## Estructuras (Structs)

### Entidades

**Patrón**: `[Descripción][Tipo]`

```go
// ✅ Correcto
type ResponseGetCurrentTime struct {
    StatusCode int       `json:"StatusCode"`
    Message    string    `json:"Message"`
    Result     time.Time `json:"CurrentTime"`
}

type ResponseHttp struct {
    Code    int    `json:"code"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

### DTOs

**Patrón**: `[Nombre]Dto`

```go
// ✅ Correcto
type MovementDto struct {
    Id int64 `json:"movimiento_id"`
}
```

### Schemas de Validación

**Patrón**: `[Nombre][Tipo]`

```go
// ✅ Correcto
type CurrentTimeBody struct {
    TimeZone string `json:"time_zone"`
    Date     string `json:"date"`
}

type CurrentTimeParams struct {
    Id int64 `uri:"id"`
}

type CurrentTimeQuery struct {
    Authorization string `form:"Authorization"`
}
```

## Paquetes

### Regla General

- **Minúsculas**
- **Una palabra** cuando sea posible
- **Descriptivos**

### Ejemplos del Proyecto

| Paquete | Ubicación | Justificación |
|---------|-----------|---------------|
| `entities` | `domain/entities/` | Contiene entidades |
| `use_case` | `application/use_case/` | Contiene casos de uso |
| `service_application` | `application/service_application/` | Contiene servicios |
| `repositories_impl` | `infrastructure/repositories_impl/` | Implementaciones de repositorios |
| `handlers` | `presentation/handlers/` | Handlers HTTP |
| `middleware` | `presentation/middleware/` | Middlewares |
| `dtos` | `presentation/dtos/` | Data Transfer Objects |
| `schema` | `presentation/schema/` | Schemas de validación |
| `container` | `presentation/container/` | Container de DI |
| `boostrap` | `presentation/boostrap/` | Bootstrap del servidor |

**Nota**: `boostrap` tiene un typo (debería ser `bootstrap`), pero se mantiene por consistencia con el código existente.

## Variables

### Variables Exportadas

**Patrón**: PascalCase

```go
// ✅ Correcto
var DatabaseConnectionToMongo infrastructure.DatabaseConnectionMongoInterface
var GetCurrentTimeRepository infrastructure.GetCurrentTimeRepository
var GetCurrentTimeUseCase use_case.GetCurrentTimeUseCase
```

### Variables No Exportadas

**Patrón**: camelCase

```go
// ✅ Correcto (si fuera necesario)
var appParams ApplicationParameters
var timeInstance timeStruct
```

### Funciones de Resolución

**Patrón**: `Resolve[Nombre]()`

```go
// ✅ Correcto
func ResolveGetCurrentTimeRepository() infrastructure.GetCurrentTimeRepository {
    // ...
}

func ResolveGetCurrentTimeUseCase() use_case.GetCurrentTimeUseCase {
    // ...
}
```

## Funciones

### Handlers

**Patrón**: `[Acción][Entidad]` o `[Verbo][Sustantivo]`

```go
// ✅ Correcto
func GetCurrentTime(ctx *gin.Context) {
    // ...
}
```

**Convención**: Nombres descriptivos que indican la acción y el recurso.

### Middlewares

**Patrón**: `[Nombre]Middleware()` o `[Nombre]()`

```go
// ✅ Correcto
func JWTMiddleware() gin.HandlerFunc {
    // ...
}

func BasicAuthMiddleware() gin.HandlerFunc {
    // ...
}

func LoggerMiddleware() gin.HandlerFunc {
    // ...
}

func RecoveryMiddleware() gin.HandlerFunc {
    // ...
}
```

### Funciones de Validación

**Patrón**: `Validate()`

```go
// ✅ Correcto
func (c CurrentTimeBody) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.TimeZone, validation.Required),
        validation.Field(&c.Date, validation.Date("2006-01-02"), validation.Required),
    )
}
```

## Constantes y Value Objects

### Constantes

**Patrón**: UPPER_SNAKE_CASE o PascalCase según contexto

```go
// ✅ Correcto
const DB_CON string = "postgresql://..."
const HOST_PORT string = ":3001"
const currentEnvironment = "ws.desarrollo.genexis.com.co"

var ID_POS int64 = 0
```

## Campos de Struct

### Campos Exportados

**Patrón**: PascalCase

```go
type ResponseGetCurrentTime struct {
    StatusCode int       // ✅ Exportado
    Message    string    // ✅ Exportado
    Result     time.Time // ✅ Exportado
}
```

### Campos No Exportados

**Patrón**: camelCase (raro en el proyecto, pero posible)

```go
type timeStruct struct {
    CurrentTime time.Time // No exportado (struct interno)
}
```

## Aliases de Import

### Cuando Usar Aliases

**Uso actual**: Aliases numéricos cuando hay conflictos

```go
import (
    use_case2 "genesis/pos/application/use_case"
    service_application2 "genesis/pos/domain/ports/application/service_application"
)
```

**Recomendación**: Usar aliases descriptivos cuando sea posible

```go
import (
    appUseCase "genesis/pos/application/use_case"
    domainService "genesis/pos/domain/ports/application/service_application"
)
```

## Resumen de Reglas

### ✅ DEBE

1. Interfaces sin sufijo (excepto ambigüedad)
2. Implementaciones con `Impl`
3. Repositorios con DB: `With[DB]Impl`
4. Métodos: `ExecuteUseCase()`, `ExecuteImpl()`, `ExecuteService()`
5. Paquetes en minúsculas, una palabra
6. Structs en PascalCase
7. Variables exportadas en PascalCase
8. Funciones de resolución: `Resolve[Nombre]()`

### ❌ NO DEBE

1. Interfaces con sufijo `Interface` (excepto ambigüedad)
2. Implementaciones sin `Impl`
3. Paquetes en PascalCase
4. Variables exportadas en camelCase
5. Mezclar convenciones en el mismo tipo de elemento

## Ejemplos Completos

### Caso de Uso Completo

```go
// domain/ports/application/use_case/get_user_use_case.go
type GetUserUseCase interface {
    ExecuteUseCase() (entities.User, error)
}

// application/use_case/get_user_use_case.go
type GetUserUseCaseImpl struct {
    GetUserRepo infrastructure.GetUserRepository
}

func (g *GetUserUseCaseImpl) ExecuteUseCase() (entities.User, error) {
    user, err := g.GetUserRepo.ExecuteImpl()
    return user, err
}
```

### Repositorio Completo

```go
// domain/ports/infrastructure/get_user_repository.go
type GetUserRepository interface {
    ExecuteImpl() (entities.User, error)
}

// infrastructure/repositories_impl/get_user_with_postgres.go
type GetUserWithPostgresImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}

func (r *GetUserWithPostgresImpl) ExecuteImpl() (entities.User, error) {
    // Implementación
}
```

