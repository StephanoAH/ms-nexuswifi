# Manejo de Errores

## Introducción

Este documento describe los patrones de manejo de errores detectados en el proyecto y las convenciones a seguir.

## Principio Fundamental

**Regla**: Siempre retornar `error` como último valor de retorno en funciones que pueden fallar.

```go
// ✅ Correcto
func ExecuteUseCase() (time.Time, error) {
    // ...
}

// ❌ Incorrecto
func ExecuteUseCase() time.Time {
    // ...
}
```

## Propagación de Errores

### Patrón Básico

Los errores se propagan hacia arriba en la pila de llamadas, sin modificación en la mayoría de los casos.

```go
// Repositorio
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    result, err := r.DB.Query(...)
    if err != nil {
        return time.Time{}, err  // Propagación directa
    }
    return result, nil
}

// Caso de Uso
func (uc *UseCaseImpl) ExecuteUseCase() (time.Time, error) {
    result, err := uc.Repo.ExecuteImpl()
    if err != nil {
        return time.Time{}, err  // Propagación sin modificar
    }
    return result, nil
}
```

### Propagación con Contexto

Cuando sea necesario agregar contexto, usar `fmt.Errorf` con `%w`:

```go
func (r *RepositoryImpl) ExecuteImpl() (time.Time, error) {
    result, err := r.DB.Query(...)
    if err != nil {
        return time.Time{}, fmt.Errorf("error querying current time: %w", err)
    }
    return result, nil
}
```

**Nota**: Este patrón no se usa actualmente en el proyecto, pero es una buena práctica para el futuro.

## Logging de Errores

### Patrón Actual

**Regla**: Usar `log.Println()` para logging de errores antes de retornarlos.

### En Repositorios

```go
func (getTime *GetCurrentTimeWithPostgresImpl) ExecuteImpl() (time.Time, error) {
    query := `SELECT NOW() AS "CurrentTime";`
    conn := getTime.DB.GetDatabaseConnectionPgx().UsePgxDatabaseConnectionDriver()
    resultCurrentTime := conn.QueryRow(context.Background(), query)
    
    err := resultCurrentTime.Scan(&timeInstance.CurrentTime)
    if err != nil {
        log.Println("Error ::: ", err)  // ✅ Logging antes de retornar
        return time.Time{}, err
    }
    
    return timeInstance.CurrentTime, nil
}
```

### En Casos de Uso

```go
func (g *GetSalesUseCaseImpl) ExecuteUseCase() (string, error) {
    r, err := g.GetSalesUseCaseRepo.ExecuteImpl()
    if err != nil {
        log.Println(err)  // ✅ Logging opcional
    }
    return r, nil
}
```

### En Servicios

```go
func (g *GetCurrentTimeServiceImpl) ExecuteService() entities.ResponseGetCurrentTime {
    getCurrentTime, err := g.GetCurrentTimeServiceUseCase.ExecuteUseCase()
    
    response := &entities.ResponseGetCurrentTime{}
    
    if err != nil {
        response = &entities.ResponseGetCurrentTime{
            StatusCode: 500,
            Message:    "Error querying the current date",
            Result:     getCurrentTime,
        }
        log.Println("Error ::: ", err)  // ✅ Logging antes de retornar respuesta
        return *response
    }
    
    // ...
}
```

## Manejo de Errores en Handlers

### Patrón Actual

Los handlers construyen respuestas HTTP basadas en el resultado del servicio.

```go
func GetCurrentTime(ctx *gin.Context) {
    response := container.GetCurrentTimeService.ExecuteService()
    res := &entities.ResponseHttp{}
    
    res.Code = response.StatusCode
    res.Message = response.Result.String()
    
    if response.StatusCode != 200 {
        res.Success = false
        ctx.JSON(http.StatusBadRequest, res)  // ✅ Respuesta de error
        return
    }
    
    res.Success = true
    ctx.JSON(http.StatusOK, res)
    return
}
```

### Estructura de Respuestas de Error

El proyecto usa dos estructuras principales para respuestas:

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

**Uso en Middlewares**:
```go
c.JSON(http.StatusUnauthorized, entities.ResponseJson{
    Status:  http.StatusUnauthorized,
    Message: "No Authorization header found",
    Data:    nil,
})
```

## Manejo de Errores en Middlewares

### Validación

Los middlewares de validación retornan errores HTTP directamente:

```go
func ValidateBody(obj Validatable) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(obj); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "JSON invalid",
                "details": err.Error(),
                "success": false,
                "status": http.StatusBadRequest,
            })
            c.Abort()
            return
        }
        if err := obj.Validate(); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Failed validation",
                "details": err.Error(),
                "success": false,
                "status": http.StatusBadRequest,
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### Autenticación

Los middlewares de autenticación retornan errores HTTP estructurados:

```go
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.Request.Header.Get("Authorization")
        
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, entities.ResponseJson{
                Status:  http.StatusBadRequest,
                Message: "No Authorization header found",
                Data:    nil,
            })
            c.Abort()
            return
        }
        
        // ...
    }
}
```

## Manejo de Panics

### Recovery Middleware

El proyecto incluye un middleware de recovery para manejar panics:

```go
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("Panic recovered: %v\n", err)
                
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal Server Error",
                })
                
                c.Abort()
            }
        }()
        
        c.Next()
    }
}
```

**Características**:
- Captura panics con `recover()`
- Retorna respuesta HTTP 500
- Aborta la cadena de handlers
- Logging del panic (actualmente con `fmt.Printf`)

## Cierre de Recursos y Errores

### Patrón con Defer

Cuando se cierran recursos, los errores de cierre se manejan en el defer:

```go
func (m *GetSalesWithMongoImpl) ExecuteImpl() (string, error) {
    result, err := collection.Find(context.Background(), bson.D{})
    if err != nil {
        log.Println("Error ::: ", err)
        return "", err
    }
    
    defer func(result *mongo.Cursor, ctx context.Context) {
        err := result.Close(ctx)
        if err != nil {
            // Manejo silencioso o logging opcional
        }
    }(result, context.Background())
    
    // ...
}
```

**Nota**: El proyecto actualmente maneja errores de cierre de forma silenciosa. Para producción, considerar logging.

## Errores en Conexiones de Base de Datos

### Patrón Actual

Las conexiones de base de datos usan `log.Fatal` o `log.Println`:

```go
func (configConnection *ConfigConnectionPgx) GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver {
    connection, err := pgx.Connect(context.Background(), configConnection.UrlToConnect)
    if err != nil {
        log.Fatal(err)  // ⚠️ Termina la aplicación
    }
    return &postgres.ConfigDatabaseConnectionPgxDriver{PgxConn: connection}
}
```

**Consideración**: `log.Fatal` termina la aplicación. Para producción, considerar retornar error y manejar en el bootstrap.

## Convenciones por Capa

### Domain

- **NO maneja errores**: Solo define estructuras
- Las interfaces de puertos definen `error` en sus firmas

### Application

- **Propaga errores**: No modifica, solo propaga
- **Logging opcional**: En casos de uso y servicios
- **Transforma a entidades de respuesta**: En servicios cuando hay error

### Infrastructure

- **Logging antes de retornar**: `log.Println("Error ::: ", err)`
- **Propaga errores**: Retorna error original
- **Manejo de recursos**: Defer para cierre, errores de cierre opcionales

### Presentation

- **Transforma errores a HTTP**: Códigos de estado apropiados
- **Respuestas estructuradas**: Usa `ResponseHttp` o `ResponseJson`
- **Validación temprana**: Middlewares validan antes de llegar a handlers

## Mejores Prácticas Recomendadas

### 1. Errores Tipados

Para el futuro, considerar errores tipados:

```go
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

### 2. Contexto con fmt.Errorf

Agregar contexto cuando se propaga:

```go
if err != nil {
    return fmt.Errorf("failed to get current time: %w", err)
}
```

### 3. Logging Estructurado

Considerar logging estructurado para producción:

```go
log.Printf("error=%v component=repository operation=get_current_time", err)
```

### 4. Retornar Errores en Conexiones

En lugar de `log.Fatal`, retornar error:

```go
func (c *ConfigConnectionPgx) GetDatabaseConnectionPgx() (*postgres.ConfigDatabaseConnectionPgxDriver, error) {
    connection, err := pgx.Connect(context.Background(), c.UrlToConnect)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    return &postgres.ConfigDatabaseConnectionPgxDriver{PgxConn: connection}, nil
}
```

## Resumen

### ✅ Patrones Actuales

1. Siempre retornar `error` como último valor
2. Usar `log.Println()` para logging de errores
3. Propagación directa de errores (sin modificación)
4. Transformación a respuestas HTTP en handlers
5. Middlewares retornan errores HTTP directamente
6. Recovery middleware para panics
7. Defer para cierre de recursos

### 📝 Recomendaciones Futuras

1. Errores tipados para mejor manejo
2. Contexto con `fmt.Errorf` y `%w`
3. Logging estructurado
4. Retornar errores en lugar de `log.Fatal`
5. Manejo explícito de errores de cierre de recursos

