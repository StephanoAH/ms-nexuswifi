# Estructura Hexagonal Detallada

## Concepto de Arquitectura Hexagonal

La Arquitectura Hexagonal (Ports and Adapters) organiza el código en un **hexágono** donde:

- **El centro (núcleo)**: Contiene la lógica de negocio pura
- **Los puertos**: Interfaces que definen contratos
- **Los adaptadores**: Implementaciones concretas que conectan con el mundo exterior

## El Hexágono en Nuestro Proyecto

```
                    ┌─────────────────┐
                    │   Presentation  │
                    │   (Adaptadores   │
                    │    de Entrada)   │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   Application   │
                    │   (Casos de     │
                    │     Uso)        │
                    └────────┬────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
┌───────▼────────┐  ┌────────▼────────┐  ┌───────▼────────┐
│    Domain      │  │   Domain        │  │    Domain      │
│   Entities     │  │   Ports         │  │  Value Objects │
│                │  │  (Interfaces)   │  │                │
└────────────────┘  └────────┬────────┘  └────────────────┘
                             │
                    ┌────────▼────────┐
                    │ Infrastructure  │
                    │  (Adaptadores   │
                    │   de Salida)    │
                    └─────────────────┘
```

## Capas Detalladas

### 1. Domain Layer (Núcleo)

#### 1.1 Entities (`domain/entities/`)

Las entidades representan conceptos del dominio de negocio. Son estructuras de datos que encapsulan lógica de negocio.

**Características**:
- Son estructuras Go simples
- Pueden tener métodos asociados
- Representan conceptos del negocio
- **NO tienen dependencias externas**

**Ejemplos en el proyecto**:
- `ResponseGetCurrentTime`: Respuesta estructurada con tiempo actual
- `BasicAuth`: Credenciales de autenticación
- `VentaDetallada`: Entidad de venta desde MongoDB
- `ResponseHttp`: Respuesta HTTP genérica
- `ResponseJson`: Respuesta JSON estructurada

**Patrón de Naming**:
- Nombres descriptivos en PascalCase
- Reflejan conceptos del dominio

#### 1.2 Ports (`domain/ports/`)

Los puertos son **interfaces** que definen contratos. Se dividen en:

**a) Application Ports** (`domain/ports/application/`):
- Definen contratos para casos de uso
- Definen contratos para servicios de aplicación
- Ejemplo: `GetCurrentTimeUseCase`, `GetCurrentTimeService`

**b) Infrastructure Ports** (`domain/ports/infrastructure/`):
- Definen contratos para repositorios
- Definen contratos para conexiones de base de datos
- Ejemplo: `GetCurrentTimeRepository`, `DatabaseConnectionPgxInterface`

**Características de los Puertos**:
- Son interfaces Go puras
- No tienen implementación
- Definen **QUÉ** se necesita, no **CÓMO** se implementa
- Permiten inversión de dependencias

**Ejemplo de Puerto de Aplicación**:
```go
type GetCurrentTimeUseCase interface {
    ExecuteUseCase() (time.Time, error)
}
```

**Ejemplo de Puerto de Infraestructura**:
```go
type GetCurrentTimeRepository interface {
    ExecuteImpl() (time.Time, error)
}
```

#### 1.3 Value Objects (`domain/value_object/`)

Los value objects son objetos inmutables que representan conceptos del dominio.

**Estructura actual**:
- `constants/`: Constantes y valores constantes del dominio
  - `running_server.go`: Constantes del servidor
  - `url_service.go`: URLs de servicios

### 2. Application Layer

#### 2.1 Use Cases (`application/use_case/`)

Los casos de uso implementan la lógica de aplicación. Cada caso de uso:

- Implementa un puerto de `domain/ports/application/use_case/`
- Orquesta llamadas a repositorios
- Contiene lógica de negocio específica de la aplicación
- **NO contiene lógica de infraestructura**

**Patrón de Implementación**:
```go
type GetCurrentTimeUseCaseImpl struct {
    GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository
}

func (g *GetCurrentTimeUseCaseImpl) ExecuteUseCase() (time.Time, error) {
    CurrentTime, err := g.GetCurrentTimeRepo.ExecuteImpl()
    return CurrentTime, err
}
```

**Características**:
- Estructura con sufijo `Impl`
- Recibe dependencias a través del constructor (struct fields)
- Método principal: `ExecuteUseCase()`
- Retorna datos del dominio o errores

#### 2.2 Services (`application/service_application/`)

Los servicios de aplicación orquestan múltiples casos de uso y coordinan operaciones complejas.

**Características**:
- Implementan puertos de `domain/ports/application/service_application/`
- Pueden llamar a múltiples casos de uso
- Transforman respuestas de casos de uso en entidades de respuesta
- Manejan lógica de orquestación

**Ejemplo**:
```go
type GetCurrentTimeServiceImpl struct {
    GetCurrentTimeServiceUseCase use_case.GetCurrentTimeUseCase
    GetSalesUseCase              use_case.GetSalesUseCase
}

func (g *GetCurrentTimeServiceImpl) ExecuteService() entities.ResponseGetCurrentTime {
    // Orquesta múltiples casos de uso
    // Transforma respuestas
    // Maneja errores
}
```

### 3. Infrastructure Layer

#### 3.1 Repository Implementations (`infrastructure/repositories_impl/`)

Las implementaciones de repositorios proporcionan acceso a datos concretos.

**Características**:
- Implementan puertos de `domain/ports/infrastructure/`
- Contienen lógica específica de acceso a datos
- Manejan conexiones a bases de datos
- Ejecutan queries SQL o operaciones NoSQL

**Patrón de Naming**:
- `GetXxxWithPostgresImpl`: Implementación con PostgreSQL
- `GetXxxWithMongoImpl`: Implementación con MongoDB
- `XxxRepositoryImpl`: Implementación genérica

**Ejemplo**:
```go
type GetCurrentTimeWithPostgresImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}

func (getTime *GetCurrentTimeWithPostgresImpl) ExecuteImpl() (time.Time, error) {
    // Lógica de acceso a PostgreSQL
}
```

#### 3.2 Database Drivers (`infrastructure/db/`)

Los drivers encapsulan la configuración y conexión a bases de datos.

**Estructura**:
- `postgres/drivers/`: Configuración PostgreSQL con pgx
- `mongo/drivers/`: Configuración MongoDB con mongo-driver

**Patrón**:
- Cada driver implementa una interfaz de conexión
- Proporciona métodos para obtener conexiones activas
- Maneja configuración de conexión

#### 3.3 External HTTP (`infrastructure/external_http/`)

Cliente HTTP para comunicarse con APIs externas.

**Estructura**:
- `http_client_interface.go`: Interfaz del cliente
- `net_http/`: Implementación con `net/http` estándar

### 4. Presentation Layer

#### 4.1 Handlers (`presentation/handlers/`)

Los handlers son funciones que procesan requests HTTP.

**Características**:
- Reciben `*gin.Context`
- Validan y transforman requests
- Llaman a servicios del container
- Construyen respuestas HTTP

**Patrón**:
```go
func GetCurrentTime(ctx *gin.Context) {
    response := container.GetCurrentTimeService.ExecuteService()
    // Construir respuesta HTTP
    ctx.JSON(http.StatusOK, res)
}
```

#### 4.2 Routers (`presentation/routers/`)

Configuran las rutas HTTP y aplican middlewares.

**Patrón**:
```go
func CurrentTime(router *gin.RouterGroup) {
    router.POST("/:id",
        middleware.ValidateParams(...),
        middleware.ValidateQuery(...),
        middleware.ValidateBody(...),
        handlers.GetCurrentTime)
}
```

#### 4.3 Middlewares (`presentation/middleware/`)

Middlewares para cross-cutting concerns:

- **JWT Middleware**: Validación de tokens JWT
- **Basic Auth Middleware**: Autenticación básica
- **Logger Middleware**: Logging de requests/responses
- **Recovery Middleware**: Manejo de panics
- **Validation Middlewares**: Validación de params, query, body

#### 4.4 Container (`presentation/container/`)

Implementa inyección de dependencias manual (patrón Container).

**Características**:
- Variables globales para instancias
- Funciones `ResolveXxx()` para obtener instancias
- Patrón Singleton (si ya existe, retorna la misma)
- Inicializa todas las dependencias en `StartContainer()`

**Flujo**:
1. Resuelve conexiones de base de datos
2. Resuelve repositorios
3. Resuelve casos de uso
4. Resuelve servicios

#### 4.5 DTOs y Schemas

- **DTOs** (`presentation/dtos/`): Data Transfer Objects para requests
- **Schemas** (`presentation/schema/`): Schemas con validación usando `ozzo-validation`

## Flujo Completo de un Request

```
1. Request HTTP llega
   ↓
2. Bootstrap carga middlewares globales
   ↓
3. Router enruta a handler específico
   ↓
4. Middlewares de ruta validan (params, query, body)
   ↓
5. Handler obtiene servicio del container
   ↓
6. Service ejecuta casos de uso
   ↓
7. Use Case llama a repositorio (a través de puerto)
   ↓
8. Repository Implementation accede a base de datos
   ↓
9. Respuesta fluye de vuelta por las capas
   ↓
10. Handler construye respuesta HTTP
```

## Reglas de la Arquitectura Hexagonal

1. **El dominio es independiente**: No conoce nada de infraestructura o presentación
2. **Las dependencias apuntan hacia adentro**: Capas externas dependen de internas
3. **Los puertos definen contratos**: Interfaces en el dominio
4. **Los adaptadores implementan**: Infraestructura y presentación implementan puertos
5. **La aplicación orquesta**: Coordina entre dominio e infraestructura

## Ventajas Específicas

- **Cambiar de base de datos**: Solo cambiar implementación en `infrastructure/repositories_impl/`
- **Cambiar de framework HTTP**: Solo cambiar `presentation/handlers/` y `presentation/routers/`
- **Testing**: Mockear puertos es trivial
- **Múltiples interfaces**: Fácil agregar CLI, gRPC, etc. sin tocar lógica de negocio

