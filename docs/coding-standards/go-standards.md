# Estándares de Go

## Introducción

Este documento define los estándares de código Go específicos del proyecto, basados en las convenciones y patrones detectados en el código existente.

## Versión de Go

El proyecto usa **Go 1.22.4** (definido en `go.mod`).

## Estructura de Paquetes

### Organización

Cada directorio representa un paquete. Los nombres de paquetes deben ser:

- **En minúsculas**
- **Una sola palabra** cuando sea posible
- **Descriptivos** del contenido

**Ejemplos del proyecto**:
- `entities` (no `Entities` o `entity`)
- `use_case` (no `usecase` o `UseCase`)
- `repositories_impl` (no `repositories` o `RepositoriesImpl`)

### Naming de Paquetes

| Tipo | Convención | Ejemplo |
|------|-----------|---------|
| Entidades | `entities` | `domain/entities/` |
| Casos de uso | `use_case` | `application/use_case/` |
| Servicios | `service_application` | `application/service_application/` |
| Repositorios | `repositories_impl` | `infrastructure/repositories_impl/` |
| Handlers | `handlers` | `presentation/handlers/` |
| Middlewares | `middleware` | `presentation/middleware/` |
| DTOs | `dtos` | `presentation/dtos/` |
| Schemas | `schema` | `presentation/schema/` |

## Convenciones de Nombres

### Interfaces

**Regla**: Interfaces **NO llevan sufijo** (a menos que haya ambigüedad).

**Ejemplos correctos**:
```go
type GetCurrentTimeUseCase interface {
    ExecuteUseCase() (time.Time, error)
}

type GetCurrentTimeRepository interface {
    ExecuteImpl() (time.Time, error)
}
```

**Excepciones** (cuando hay ambigüedad):
```go
type BasicAuthRepositoryPort interface {  // "Port" para claridad
    GetCredentials(...) (entities.BasicAuth, error)
}

type DatabaseConnectionPgxInterface interface {  // "Interface" para claridad
    GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver
}
```

### Implementaciones

**Regla**: Implementaciones llevan sufijo `Impl`.

**Patrón**: `[Nombre]Impl`

**Ejemplos**:
```go
type GetCurrentTimeUseCaseImpl struct {
    GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository
}

type GetCurrentTimeWithPostgresImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}

type BasicAuthRepositoryImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}
```

### Métodos

#### Casos de Uso

**Patrón**: `ExecuteUseCase()`

```go
func (g *GetCurrentTimeUseCaseImpl) ExecuteUseCase() (time.Time, error) {
    CurrentTime, err := g.GetCurrentTimeRepo.ExecuteImpl()
    return CurrentTime, err
}
```

#### Repositorios

**Patrón**: `ExecuteImpl()` o métodos específicos del dominio

```go
// Patrón genérico
func (r *RepositoryImpl) ExecuteImpl() (T, error) {
    // ...
}

// Patrón específico del dominio
func (b *BasicAuthRepositoryImpl) GetCredentials(integrationID int64, username string) (entities.BasicAuth, error) {
    // ...
}
```

#### Servicios

**Patrón**: `ExecuteService()`

```go
func (g *GetCurrentTimeServiceImpl) ExecuteService() entities.ResponseGetCurrentTime {
    // ...
}
```

### Estructuras

**Regla**: PascalCase, nombres descriptivos.

**Ejemplos**:
```go
type ResponseGetCurrentTime struct {
    StatusCode int       `json:"StatusCode"`
    Message    string    `json:"Message"`
    Result     time.Time `json:"CurrentTime"`
}

type GetCurrentTimeUseCaseImpl struct {
    GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository
}
```

### Variables y Campos

**Regla**: 
- **Exportadas**: PascalCase
- **No exportadas**: camelCase

**Ejemplos**:
```go
type Response struct {
    StatusCode int    // Exportado
    message    string // No exportado (si fuera necesario)
}

var DatabaseConnectionToMongo infrastructure.DatabaseConnectionMongoInterface // Exportado
var appParams ApplicationParameters // No exportado (si fuera necesario)
```

## Manejo de Errores

### Convención General

**Regla**: Siempre retornar `error` como último valor de retorno.

```go
func ExecuteUseCase() (time.Time, error) {
    // ...
}
```

### Propagación de Errores

**Patrón detectado**: Los errores se propagan hacia arriba, con logging opcional.

```go
// En repositorio
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    result, err := r.DB.Query(...)
    if err != nil {
        log.Println("Error ::: ", err)  // Logging opcional
        return time.Time{}, err  // Propagación
    }
    return result, nil
}

// En caso de uso
func (uc *UseCaseImpl) ExecuteUseCase() (time.Time, error) {
    result, err := uc.Repo.ExecuteImpl()
    if err != nil {
        return time.Time{}, err  // Propagación sin modificar
    }
    return result, nil
}
```

### Logging de Errores

**Patrón**: Usar `log.Println()` para logging de errores.

```go
if err != nil {
    log.Println("Error ::: ", err)
    return nil, err
}
```

**Ubicaciones típicas**:
- Repositorios: Logging antes de retornar error
- Servicios: Logging antes de construir respuesta de error

## Tags de Struct

### Tags JSON

**Regla**: Siempre incluir tags `json:` para estructuras que se serializan.

```go
type ResponseGetCurrentTime struct {
    StatusCode int       `json:"StatusCode"`
    Message    string    `json:"Message"`
    Result     time.Time `json:"CurrentTime"`
}
```

**Convención**: 
- Usar PascalCase en JSON tags (coincide con nombres de campos)
- O usar snake_case si el API lo requiere

### Tags DB

**Regla**: Incluir tags `db:` cuando se mapean directamente a columnas de base de datos.

```go
type BasicAuth struct {
    Username string `json:"username" db:"identificador_conexion"`
    Password string `json:"password" db:"clave_conexion"`
}
```

## Inyección de Dependencias

### Patrón Manual (Container)

El proyecto usa inyección de dependencias **manual** a través del Container.

**Patrón**:
```go
// 1. Variable global
var GetCurrentTimeRepository infrastructure.GetCurrentTimeRepository

// 2. Función de resolución (Singleton)
func ResolveGetCurrentTimeRepository() infrastructure.GetCurrentTimeRepository {
    if GetCurrentTimeRepository != nil {
        return GetCurrentTimeRepository
    }
    
    GetCurrentTimeRepository = &repositories_impl.GetCurrentTimeWithPostgresImpl{
        DB: DatabaseConnectionToLecWithPgx,
    }
    return GetCurrentTimeRepository
}

// 3. Inicialización en StartContainer()
func StartContainer() {
    GetCurrentTimeRepository = ResolveGetCurrentTimeRepository()
    // ...
}
```

**Reglas**:
- Variables globales para instancias
- Funciones `ResolveXxx()` para obtener instancias
- Patrón Singleton: si ya existe, retorna la misma
- Inicialización centralizada en `StartContainer()`

## Imports

### Organización

**Orden**:
1. Paquetes estándar de Go
2. Paquetes de terceros
3. Paquetes del proyecto

**Ejemplo**:
```go
import (
    "context"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
    
    "genesis/pos/domain/entities"
    "genesis/pos/domain/ports/infrastructure"
)
```

### Aliases

**Uso de aliases** cuando hay conflictos de nombres:

```go
import (
    use_case2 "genesis/pos/application/use_case"
    service_application2 "genesis/pos/domain/ports/application/service_application"
    drivers2 "genesis/pos/infrastructure/db/postgres/drivers"
)
```

**Nota**: En el código actual se usan aliases numéricos (`use_case2`), pero es preferible usar aliases descriptivos cuando sea posible.

## Context

### Uso de Context

**Regla**: Usar `context.Background()` para operaciones de base de datos cuando no hay contexto de request.

```go
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    ctx := context.Background()
    query := `SELECT NOW() AS "CurrentTime";`
    conn := r.DB.GetDatabaseConnectionPgx().UsePgxDatabaseConnectionDriver()
    result := conn.QueryRow(ctx, query)
    // ...
}
```

**Nota**: Para futuras mejoras, considerar pasar `context.Context` desde los handlers para propagar timeouts y cancelación.

## Cierre de Recursos

### Defer para Cierre

**Regla**: Usar `defer` para cerrar recursos (conexiones, cursors, etc.).

```go
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    conn := r.DB.GetDatabaseConnectionPgx().UsePgxDatabaseConnectionDriver()
    defer conn.Close(context.Background())
    // ...
}

func (m *GetSalesWithMongoImpl) ExecuteImpl() (string, error) {
    result, err := collection.Find(context.Background(), bson.D{})
    if err != nil {
        return "", err
    }
    defer func(result *mongo.Cursor, ctx context.Context) {
        err := result.Close(ctx)
        if err != nil {
            // Manejo de error en cierre
        }
    }(result, context.Background())
    // ...
}
```

## Comentarios

### Documentación

**Regla**: Comentarios descriptivos cuando el código no es autoexplicativo.

**Ejemplos del proyecto**:
```go
// Check if the content type is multipart/form-data
if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
    // Log request without body if it's multipart (likely contains files)
    log.Printf("Request [multipart]: %s %s from %s", ...)
}
```

### Comentarios de Código Comentado

**Evitar**: Código comentado en producción. Si es necesario para referencia, usar comentarios explicativos.

```go
// var requestBody dtos.MovementDto
// if err := ctx.BindJSON(&requestBody); err != nil {
//     ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//     return
// }
```

**Mejor**: Eliminar código comentado o documentar por qué está comentado.

## Formato

### gofmt

**Regla**: Siempre usar `gofmt` o formateo automático del IDE.

```bash
gofmt -w .
```

### Líneas de Código

**Regla**: Mantener líneas razonablemente cortas (< 120 caracteres cuando sea posible).

## Testing

### Naming de Tests

**Patrón**: `[nombre_funcion]_test.go`

**Ejemplo**: Para `get_current_time_use_case.go`, el test sería `get_current_time_use_case_test.go`.

### Estructura de Tests

Ver documentación en `/docs/testing/` para detalles completos.

## Resumen de Reglas Clave

1. ✅ Interfaces sin sufijo (excepto ambigüedad)
2. ✅ Implementaciones con sufijo `Impl`
3. ✅ Métodos: `ExecuteUseCase()`, `ExecuteImpl()`, `ExecuteService()`
4. ✅ Siempre retornar `error` como último valor
5. ✅ Usar `log.Println()` para logging de errores
6. ✅ Tags `json:` y `db:` cuando corresponda
7. ✅ Inyección de dependencias manual vía Container
8. ✅ Usar `defer` para cerrar recursos
9. ✅ Imports organizados: stdlib → terceros → proyecto

