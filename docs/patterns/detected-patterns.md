# Patrones Detectados en el Proyecto

## Introducción

Este documento lista todos los patrones arquitectónicos y de diseño detectados en el código del proyecto. Estos patrones deben mantenerse y seguirse en nuevas implementaciones.

## Patrones Arquitectónicos

### 1. Arquitectura Hexagonal (Ports and Adapters)

**Descripción**: Separación del código en capas con dependencias hacia adentro.

**Implementación**:
- `domain/`: Núcleo sin dependencias
- `application/`: Depende solo de `domain/`
- `infrastructure/`: Implementa puertos de `domain/`
- `presentation/`: Depende de `domain/` y `application/`

**Ejemplo**:
```
domain/ports/infrastructure/get_current_time_repository.go (puerto)
    ↑ implementa
infrastructure/repositories_impl/get_current_time_with_postgres.go
```

### 2. Inversión de Dependencias (Dependency Inversion)

**Descripción**: Depender de abstracciones (interfaces), no de implementaciones concretas.

**Implementación**:
- Interfaces en `domain/ports/`
- Implementaciones en `application/` e `infrastructure/`
- Casos de uso dependen de interfaces, no de repositorios concretos

**Ejemplo**:
```go
// Caso de uso depende de interfaz
type GetCurrentTimeUseCaseImpl struct {
    GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository  // Interfaz
}
```

### 3. Container Pattern (Inyección de Dependencias Manual)

**Descripción**: Container manual para resolver y gestionar dependencias.

**Implementación**:
- Variables globales para instancias en `presentation/container/`
- Funciones `ResolveXxx()` para obtener instancias
- Patrón Singleton: si ya existe, retorna la misma
- Inicialización centralizada en `StartContainer()`

**Ejemplo**:
```go
var GetCurrentTimeRepository infrastructure.GetCurrentTimeRepository

func ResolveGetCurrentTimeRepository() infrastructure.GetCurrentTimeRepository {
    if GetCurrentTimeRepository != nil {
        return GetCurrentTimeRepository
    }
    GetCurrentTimeRepository = &repositories_impl.GetCurrentTimeWithPostgresImpl{
        DB: DatabaseConnectionToLecWithPgx,
    }
    return GetCurrentTimeRepository
}
```

## Patrones de Diseño

### 4. Repository Pattern

**Descripción**: Abstracción del acceso a datos.

**Implementación**:
- Interfaces en `domain/ports/infrastructure/`
- Implementaciones en `infrastructure/repositories_impl/`
- Naming: `GetXxxWithPostgresImpl`, `GetXxxWithMongoImpl`

**Ejemplo**:
```go
// Puerto
type GetCurrentTimeRepository interface {
    ExecuteImpl() (time.Time, error)
}

// Implementación
type GetCurrentTimeWithPostgresImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}
```

### 5. Use Case Pattern

**Descripción**: Encapsula lógica de aplicación en casos de uso.

**Implementación**:
- Interfaces en `domain/ports/application/use_case/`
- Implementaciones en `application/use_case/`
- Método estándar: `ExecuteUseCase()`

**Ejemplo**:
```go
type GetCurrentTimeUseCase interface {
    ExecuteUseCase() (time.Time, error)
}

type GetCurrentTimeUseCaseImpl struct {
    GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository
}
```

### 6. Service Layer Pattern

**Descripción**: Servicios que orquestan múltiples casos de uso.

**Implementación**:
- Interfaces en `domain/ports/application/service_application/`
- Implementaciones en `application/service_application/`
- Método estándar: `ExecuteService()`

**Ejemplo**:
```go
type GetCurrentTimeServiceImpl struct {
    GetCurrentTimeServiceUseCase use_case.GetCurrentTimeUseCase
    GetSalesUseCase              use_case.GetSalesUseCase
}
```

### 7. DTO Pattern (Data Transfer Object)

**Descripción**: Objetos para transferir datos entre capas.

**Implementación**:
- DTOs en `presentation/dtos/` para requests simples
- Schemas en `presentation/schema/` para validación
- Entidades en `domain/entities/` para lógica de negocio

**Ejemplo**:
```go
// DTO
type MovementDto struct {
    Id int64 `json:"movimiento_id"`
}

// Schema con validación
type CurrentTimeBody struct {
    TimeZone string `json:"time_zone"`
    Date     string `json:"date"`
}
```

### 8. Middleware Pattern

**Descripción**: Cross-cutting concerns como middlewares.

**Implementación**:
- Middlewares en `presentation/middleware/`
- Retornan `gin.HandlerFunc`
- Aplicados en rutas o globalmente

**Ejemplos**:
- `JWTMiddleware()`: Validación de tokens
- `BasicAuthMiddleware()`: Autenticación básica
- `LoggerMiddleware()`: Logging de requests
- `RecoveryMiddleware()`: Manejo de panics
- `ValidateBody()`, `ValidateQuery()`, `ValidateParams()`: Validación

### 9. Strategy Pattern (Validación)

**Descripción**: Diferentes estrategias de validación.

**Implementación**:
- Interfaz `Validatable` con método `Validate()`
- Implementaciones en schemas con `ozzo-validation`
- Middlewares genéricos que usan la interfaz

**Ejemplo**:
```go
type Validatable interface {
    Validate() error
}

func ValidateBody(obj Validatable) gin.HandlerFunc {
    // Usa la interfaz, no implementación específica
}
```

## Patrones de Naming

### 10. Interface Naming

**Patrón**: Interfaces sin sufijo (excepto ambigüedad).

**Ejemplos**:
- `GetCurrentTimeUseCase` (no `GetCurrentTimeUseCaseInterface`)
- `GetCurrentTimeRepository` (no `GetCurrentTimeRepositoryInterface`)
- `BasicAuthRepositoryPort` (excepción por ambigüedad)

### 11. Implementation Naming

**Patrón**: Implementaciones con sufijo `Impl`.

**Ejemplos**:
- `GetCurrentTimeUseCaseImpl`
- `GetCurrentTimeServiceImpl`
- `GetCurrentTimeWithPostgresImpl`

### 12. Method Naming

**Patrones**:
- Casos de uso: `ExecuteUseCase()`
- Repositorios: `ExecuteImpl()` o métodos específicos
- Servicios: `ExecuteService()`

## Patrones de Manejo de Errores

### 13. Error Propagation

**Descripción**: Errores se propagan hacia arriba sin modificación.

**Implementación**:
- Retornar `error` como último valor
- Propagación directa en la mayoría de casos
- Logging opcional antes de retornar

**Ejemplo**:
```go
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    result, err := r.DB.Query(...)
    if err != nil {
        log.Println("Error ::: ", err)
        return time.Time{}, err  // Propagación directa
    }
    return result, nil
}
```

### 14. Error Transformation

**Descripción**: Transformar errores a respuestas HTTP en handlers.

**Implementación**:
- Servicios retornan entidades con `StatusCode`
- Handlers construyen respuestas HTTP basadas en status
- Middlewares retornan errores HTTP directamente

**Ejemplo**:
```go
func GetCurrentTime(ctx *gin.Context) {
    response := container.GetCurrentTimeService.ExecuteService()
    if response.StatusCode != 200 {
        res.Success = false
        ctx.JSON(http.StatusBadRequest, res)
        return
    }
}
```

## Patrones de Acceso a Datos

### 15. Driver Pattern

**Descripción**: Drivers encapsulan configuración de conexión.

**Implementación**:
- Drivers en `infrastructure/db/[db]/drivers/`
- Interfaces de conexión en `domain/ports/infrastructure/`
- Métodos para obtener conexiones activas

**Ejemplo**:
```go
type ConfigConnectionPgx struct {
    UrlToConnect string
}

func (c *ConfigConnectionPgx) GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver {
    // Configuración y conexión
}
```

### 16. Resource Management

**Descripción**: Uso de `defer` para cierre de recursos.

**Implementación**:
- `defer conn.Close(context.Background())` para conexiones
- `defer result.Close(ctx)` para cursors
- Manejo de errores de cierre opcional

**Ejemplo**:
```go
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    conn := r.DB.GetDatabaseConnectionPgx().UsePgxDatabaseConnectionDriver()
    defer conn.Close(context.Background())
    // ...
}
```

## Patrones de Validación

### 17. Schema Validation

**Descripción**: Schemas con validación usando `ozzo-validation`.

**Implementación**:
- Structs en `presentation/schema/`
- Método `Validate()` que implementa `Validatable`
- Validación en middlewares antes de llegar a handlers

**Ejemplo**:
```go
type CurrentTimeBody struct {
    TimeZone string `json:"time_zone"`
    Date     string `json:"date"`
}

func (c CurrentTimeBody) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.TimeZone, validation.Required),
        validation.Field(&c.Date, validation.Date("2006-01-02"), validation.Required),
    )
}
```

## Patrones de Bootstrap

### 18. Bootstrap Pattern

**Descripción**: Inicialización centralizada de la aplicación.

**Implementación**:
- `UploadEnv()`: Carga variables de entorno
- `StartContainer()`: Inicializa dependencias
- `RunServer()`: Inicia el servidor HTTP

**Ejemplo**:
```go
func main() {
    boostrap.UploadEnv()
    container.StartContainer()
    boostrap.RunServer()
}
```

## Resumen de Patrones

| Patrón | Ubicación | Propósito |
|--------|-----------|-----------|
| **Arquitectura Hexagonal** | Toda la estructura | Separación de capas |
| **Inversión de Dependencias** | `domain/ports/` | Depender de abstracciones |
| **Container Pattern** | `presentation/container/` | Inyección de dependencias |
| **Repository Pattern** | `infrastructure/repositories_impl/` | Abstracción de datos |
| **Use Case Pattern** | `application/use_case/` | Lógica de aplicación |
| **Service Layer** | `application/service_application/` | Orquestación |
| **DTO Pattern** | `presentation/dtos/`, `presentation/schema/` | Transferencia de datos |
| **Middleware Pattern** | `presentation/middleware/` | Cross-cutting concerns |
| **Strategy (Validación)** | `presentation/schema/` | Validación flexible |
| **Driver Pattern** | `infrastructure/db/` | Configuración de conexión |
| **Resource Management** | Repositorios | Cierre de recursos |
| **Bootstrap Pattern** | `presentation/boostrap/` | Inicialización |

## Reglas para Mantener Consistencia

1. ✅ **Seguir patrones establecidos**: No inventar nuevos patrones sin justificación
2. ✅ **Mantener naming consistente**: Seguir convenciones detectadas
3. ✅ **Respetar dependencias**: No violar arquitectura hexagonal
4. ✅ **Usar interfaces**: Depender de puertos, no implementaciones
5. ✅ **Container para DI**: Usar container para resolver dependencias
6. ✅ **Validación en schemas**: Usar schemas con `ozzo-validation`
7. ✅ **Manejo de errores**: Propagación y transformación según capa

## Siguientes Pasos

- Ver [Reglas de Consistencia](./consistency-rules.md)
- Ver [Decisiones Arquitectónicas](./adr-initial.md)

