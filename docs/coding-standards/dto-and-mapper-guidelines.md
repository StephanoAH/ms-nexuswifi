# Guías de DTOs y Mappers

## Introducción

Este documento describe el uso de DTOs (Data Transfer Objects) y las convenciones para mapeo entre capas en el proyecto.

## Concepto de DTOs

Los DTOs son objetos que se usan para transferir datos entre capas, especialmente en la capa de presentación para recibir y enviar datos HTTP.

## Ubicación de DTOs

### DTOs de Request

**Ubicación**: `presentation/dtos/`

**Propósito**: Estructuras para recibir datos de requests HTTP.

**Ejemplo**:
```go
// presentation/dtos/movement_dto.go
package dtos

type MovementDto struct {
    Id int64 `json:"movimiento_id"`
}
```

### Schemas de Validación

**Ubicación**: `presentation/schema/`

**Propósito**: Estructuras con validación para params, query y body.

**Ejemplo**:
```go
// presentation/schema/current_time.go
package schema

import validation "github.com/go-ozzo/ozzo-validation"

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

## Diferencias: DTOs vs Schemas vs Entidades

### DTOs (`presentation/dtos/`)

- **Propósito**: Transferencia de datos HTTP simples
- **Validación**: Básica (binding de Gin)
- **Uso**: Requests simples sin validación compleja

### Schemas (`presentation/schema/`)

- **Propósito**: Validación estructurada de requests
- **Validación**: Completa con `ozzo-validation`
- **Uso**: Requests que requieren validación (params, query, body)
- **Característica**: Implementan `Validatable` con método `Validate()`

### Entidades (`domain/entities/`)

- **Propósito**: Conceptos del dominio de negocio
- **Validación**: Lógica de negocio (si aplica)
- **Uso**: Dentro del dominio y aplicación
- **Característica**: NO deben usarse directamente en HTTP

## Flujo de Datos

```
HTTP Request
    ↓
[Presentation] Schema/DTO (validación)
    ↓ (mapeo)
[Presentation] Handler (transforma a entidad o llama servicio)
    ↓
[Application] Service/Use Case (usa entidades del dominio)
    ↓
[Domain] Entity (lógica de negocio)
    ↓
[Infrastructure] Repository (persiste entidad)
```

## Patrones de Mapeo

### Patrón 1: Schema → Handler → Service

**Ejemplo con validación**:

```go
// presentation/schema/user.go
type CreateUserBody struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

func (c CreateUserBody) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.Username, validation.Required, validation.Length(3, 20)),
        validation.Field(&c.Email, validation.Required, validation.Email),
    )
}

// presentation/handlers/create_user.go
func CreateUser(ctx *gin.Context) {
    body, _ := ctx.Get("body").(*schema.CreateUserBody)
    
    // Mapeo a entidad del dominio
    user := entities.User{
        Username: body.Username,
        Email:    body.Email,
    }
    
    // Llamar a servicio
    response := container.CreateUserService.ExecuteService(user)
    // ...
}
```

### Patrón 2: DTO Simple → Handler

**Ejemplo sin validación compleja**:

```go
// presentation/dtos/movement_dto.go
type MovementDto struct {
    Id int64 `json:"movimiento_id"`
}

// presentation/handlers/get_movement.go
func GetMovement(ctx *gin.Context) {
    var dto dtos.MovementDto
    if err := ctx.BindJSON(&dto); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Usar dto directamente o mapear a entidad
    // ...
}
```

## Convenciones de Naming

### DTOs

**Patrón**: `[Nombre]Dto`

```go
type MovementDto struct {
    // ...
}

type UserDto struct {
    // ...
}
```

### Schemas

**Patrón**: `[Nombre][Tipo]`

**Tipos**:
- `Body`: Para request body
- `Params`: Para path parameters
- `Query`: Para query parameters

```go
type CurrentTimeBody struct {
    // ...
}

type CurrentTimeParams struct {
    // ...
}

type CurrentTimeQuery struct {
    // ...
}
```

## Validación

### Con ozzo-validation (Schemas)

**Patrón**: Método `Validate()` que implementa `Validatable`

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

### Uso en Middlewares

Los middlewares de validación usan la interfaz `Validatable`:

```go
type Validatable interface {
    Validate() error
}

func ValidateBody(obj Validatable) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(obj); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid"})
            c.Abort()
            return
        }
        if err := obj.Validate(); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Failed validation"})
            c.Abort()
            return
        }
        c.Set("body", obj)
        c.Next()
    }
}
```

## Mapeo a Entidades del Dominio

### En Handlers

Los handlers deben mapear DTOs/Schemas a entidades del dominio antes de llamar a servicios:

```go
func CreateUser(ctx *gin.Context) {
    body, _ := ctx.Get("body").(*schema.CreateUserBody)
    
    // Mapeo explícito
    user := entities.User{
        Username: body.Username,
        Email:    body.Email,
        CreatedAt: time.Now(),
    }
    
    response := container.CreateUserService.ExecuteService(user)
    // ...
}
```

### Evitar

**❌ NO pasar DTOs directamente a servicios**:

```go
// ❌ Incorrecto
func CreateUser(ctx *gin.Context) {
    var dto dtos.UserDto
    ctx.BindJSON(&dto)
    response := container.CreateUserService.ExecuteService(dto)  // ❌ DTO en servicio
}
```

**✅ Pasar entidades del dominio**:

```go
// ✅ Correcto
func CreateUser(ctx *gin.Context) {
    body, _ := ctx.Get("body").(*schema.CreateUserBody)
    user := entities.User{
        Username: body.Username,
        Email:    body.Email,
    }
    response := container.CreateUserService.ExecuteService(user)  // ✅ Entidad
}
```

## Mapeo de Respuestas

### Entidades a Respuestas HTTP

Los handlers transforman entidades del dominio en respuestas HTTP:

```go
func GetCurrentTime(ctx *gin.Context) {
    response := container.GetCurrentTimeService.ExecuteService()
    res := &entities.ResponseHttp{}
    
    res.Code = response.StatusCode
    res.Message = response.Result.String()
    res.Success = response.StatusCode == 200
    
    ctx.JSON(http.StatusOK, res)
}
```

### Estructuras de Respuesta

El proyecto usa dos estructuras principales:

#### ResponseHttp

```go
type ResponseHttp struct {
    Code    int    `json:"code"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

#### ResponseJson

```go
type ResponseJson struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

## Tags de Serialización

### JSON Tags

**Regla**: Siempre incluir tags `json:` en DTOs y entidades.

```go
type MovementDto struct {
    Id int64 `json:"movimiento_id"`  // ✅ Tag JSON
}
```

**Convención**: 
- Usar nombres descriptivos en JSON
- Puede diferir del nombre del campo Go

### URI Tags (Params)

```go
type CurrentTimeParams struct {
    Id int64 `uri:"id"`  // ✅ Tag URI para path params
}
```

### Form Tags (Query)

```go
type CurrentTimeQuery struct {
    Authorization string `form:"Authorization"`  // ✅ Tag form para query params
}
```

### DB Tags (Solo en Entidades)

```go
type BasicAuth struct {
    Username string `json:"username" db:"identificador_conexion"`  // ✅ Tag DB
    Password string `json:"password" db:"clave_conexion"`
}
```

## Ejemplo Completo

### 1. Schema de Validación

```go
// presentation/schema/create_user.go
package schema

import validation "github.com/go-ozzo/ozzo-validation"

type CreateUserBody struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

func (c CreateUserBody) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.Username, validation.Required, validation.Length(3, 20)),
        validation.Field(&c.Email, validation.Required, validation.Email),
        validation.Field(&c.Age, validation.Required, validation.Min(18)),
    )
}
```

### 2. Entidad del Dominio

```go
// domain/entities/user.go
package entities

import "time"

type User struct {
    ID        int64     `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"email" db:"email"`
    Age       int       `json:"age" db:"age"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### 3. Handler con Mapeo

```go
// presentation/handlers/create_user.go
package handlers

import (
    "genesis/pos/domain/entities"
    "genesis/pos/presentation/container"
    "genesis/pos/presentation/schema"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

func CreateUser(ctx *gin.Context) {
    body, exists := ctx.Get("body").(*schema.CreateUserBody)
    if !exists {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
        return
    }
    
    // Mapeo a entidad del dominio
    user := entities.User{
        Username:  body.Username,
        Email:     body.Email,
        Age:       body.Age,
        CreatedAt: time.Now(),
    }
    
    // Llamar a servicio
    response := container.CreateUserService.ExecuteService(user)
    
    // Construir respuesta HTTP
    res := &entities.ResponseHttp{
        Code:    response.StatusCode,
        Success: response.StatusCode == 200,
        Message: response.Message,
    }
    
    if response.StatusCode != 200 {
        ctx.JSON(http.StatusBadRequest, res)
        return
    }
    
    ctx.JSON(http.StatusOK, res)
}
```

### 4. Router con Validación

```go
// presentation/routers/user.go
func User(router *gin.RouterGroup) {
    router.POST("/",
        middleware.ValidateBody(new(schema.CreateUserBody)),
        handlers.CreateUser,
    )
}
```

## Reglas de Mapeo

### ✅ DEBE

1. **Mapear DTOs/Schemas a entidades** antes de llamar a servicios
2. **Usar schemas con validación** para requests complejos
3. **Usar DTOs simples** solo para requests muy simples
4. **Transformar entidades a respuestas HTTP** en handlers
5. **Incluir tags JSON** en todas las estructuras serializables

### ❌ NO DEBE

1. **Pasar DTOs a servicios** (usar entidades)
2. **Usar entidades directamente en HTTP** (usar DTOs/Schemas)
3. **Mezclar lógica de negocio en DTOs** (solo transferencia de datos)
4. **Omitir validación** en requests importantes
5. **Exponer detalles de implementación** en DTOs

## Resumen

- **DTOs**: Para transferencia simple de datos HTTP
- **Schemas**: Para validación estructurada de requests
- **Entidades**: Para lógica de negocio en el dominio
- **Mapeo**: DTOs/Schemas → Entidades → Servicios → Respuestas HTTP
- **Validación**: Usar `ozzo-validation` en schemas
- **Tags**: JSON para serialización, URI/form para binding

