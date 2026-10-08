# Capas y Dependencias

## Regla Fundamental: Dependencias Hacia Adentro

En arquitectura hexagonal, **las dependencias siempre apuntan hacia el centro** (el dominio). Las capas externas dependen de las internas, nunca al revés.

```
Presentation → Application → Domain ← Infrastructure
```

## Mapa de Dependencias

```
┌─────────────────────────────────────────────────┐
│           Presentation Layer                    │
│  (handlers, routers, middlewares, container)   │
└───────────────┬─────────────────────────────────┘
                │ depende de
┌───────────────▼─────────────────────────────────┐
│          Application Layer                      │
│     (use cases, services)                       │
└───────────────┬─────────────────────────────────┘
                │ depende de
┌───────────────▼─────────────────────────────────┐
│            Domain Layer                         │
│  (entities, ports, value objects)               │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │         NO DEPENDE DE NADIE              │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
                ↑ implementa
┌───────────────┴─────────────────────────────────┐
│         Infrastructure Layer                    │
│  (repositories, db drivers, external http)     │
└─────────────────────────────────────────────────┘
```

## Detalle por Capa

### Domain Layer

**Dependencias**: ❌ **NINGUNA**

El dominio es completamente independiente. No puede importar:

- ❌ `application/`
- ❌ `infrastructure/`
- ❌ `presentation/`
- ❌ Frameworks externos (Gin, pgx, mongo-driver)
- ❌ Paquetes de HTTP
- ❌ Paquetes de base de datos

**Puede importar**:

- ✅ Solo paquetes estándar de Go (`time`, `errors`, `context`, etc.)
- ✅ Librerías puras sin efectos secundarios

**Ejemplo Correcto**:
```go
// domain/entities/response.go
package entities

import "time"  // ✅ OK: paquete estándar

type Response struct {
    Timestamp time.Time
}
```

**Ejemplo Incorrecto**:
```go
// domain/entities/response.go
package entities

import "github.com/gin-gonic/gin"  // ❌ ERROR: framework

type Response struct {
    Context *gin.Context  // ❌ ERROR: dependencia externa
}
```

### Application Layer

**Dependencias**: ✅ Solo `domain/`

La capa de aplicación puede importar:

- ✅ `domain/` (entidades, puertos, value objects)
- ✅ Paquetes estándar de Go

**NO puede importar**:

- ❌ `infrastructure/` (usa puertos, no implementaciones)
- ❌ `presentation/` (no conoce detalles de HTTP)

**Ejemplo Correcto**:
```go
// application/use_case/get_time_use_case.go
package use_case

import (
    "genesis/pos/domain/ports/infrastructure"  // ✅ OK: puerto
    "time"
)

type GetTimeUseCaseImpl struct {
    Repo infrastructure.GetTimeRepository  // ✅ OK: interfaz, no implementación
}
```

**Ejemplo Incorrecto**:
```go
// application/use_case/get_time_use_case.go
package use_case

import "genesis/pos/infrastructure/repositories_impl"  // ❌ ERROR: implementación

type GetTimeUseCaseImpl struct {
    Repo *repositories_impl.GetTimeWithPostgresImpl  // ❌ ERROR: dependencia directa
}
```

### Infrastructure Layer

**Dependencias**: ✅ Solo `domain/`

La infraestructura implementa puertos del dominio, por lo que:

- ✅ Importa `domain/ports/infrastructure/` (para implementar)
- ✅ Puede importar frameworks (pgx, mongo-driver, etc.)
- ✅ Puede importar paquetes estándar

**NO puede importar**:

- ❌ `application/` (no conoce casos de uso)
- ❌ `presentation/` (no conoce HTTP)

**Ejemplo Correcto**:
```go
// infrastructure/repositories_impl/get_time_postgres.go
package repositories_impl

import (
    "genesis/pos/domain/ports/infrastructure"  // ✅ OK: puerto a implementar
    "github.com/jackc/pgx/v5"  // ✅ OK: framework de DB
)

type GetTimeWithPostgresImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface  // ✅ OK: puerto
}
```

### Presentation Layer

**Dependencias**: ✅ `domain/` y `application/`

La presentación es la capa más externa y puede:

- ✅ Importar `domain/` (para entidades y DTOs)
- ✅ Importar `application/` (para servicios y casos de uso)
- ✅ Importar frameworks HTTP (Gin)
- ✅ Importar paquetes estándar

**NO puede importar directamente**:

- ❌ `infrastructure/` (usa el container para resolver dependencias)

**Ejemplo Correcto**:
```go
// presentation/handlers/get_time.go
package handlers

import (
    "genesis/pos/domain/entities"  // ✅ OK: entidades
    "genesis/pos/presentation/container"  // ✅ OK: container resuelve dependencias
    "github.com/gin-gonic/gin"  // ✅ OK: framework HTTP
)

func GetTime(ctx *gin.Context) {
    service := container.GetTimeService  // ✅ OK: resuelto por container
    // ...
}
```

## Flujo de Resolución de Dependencias

El **Container** (`presentation/container/`) es responsable de resolver todas las dependencias:

```
1. Container crea conexiones de DB (infrastructure)
   ↓
2. Container crea repositorios (infrastructure) inyectando conexiones
   ↓
3. Container crea casos de uso (application) inyectando repositorios
   ↓
4. Container crea servicios (application) inyectando casos de uso
   ↓
5. Handlers (presentation) usan servicios del container
```

**Ejemplo del Container**:
```go
// presentation/container/container.go

// 1. Resolver conexión
var DatabaseConnectionToLecWithPgx infrastructure.DatabaseConnectionPgxInterface

func ResolveDatabaseConnectionToLecWithPgx() infrastructure.DatabaseConnectionPgxInterface {
    return &drivers2.ConfigConnectionPgx{
        UrlToConnect: os.Getenv("CONNECTION_POSTGRESQL"),
    }
}

// 2. Resolver repositorio (usa conexión)
var GetCurrentTimeRepository infrastructure.GetCurrentTimeRepository

func ResolveGetCurrentTimeRepository() infrastructure.GetCurrentTimeRepository {
    return &repositories_impl.GetCurrentTimeWithPostgresImpl{
        DB: DatabaseConnectionToLecWithPgx,  // ✅ Inyección
    }
}

// 3. Resolver caso de uso (usa repositorio)
var GetCurrentTimeUseCase use_case.GetCurrentTimeUseCase

func ResolveGetCurrentTimeUseCase() use_case.GetCurrentTimeUseCase {
    return &use_case2.GetCurrentTimeUseCaseImpl{
        GetCurrentTimeRepo: GetCurrentTimeRepository,  // ✅ Inyección
    }
}

// 4. Resolver servicio (usa caso de uso)
var GetCurrentTimeService service_application2.GetCurrentTimeService

func ResolveGetCurrentTimeService() service_application2.GetCurrentTimeService {
    return &service_application.GetCurrentTimeServiceImpl{
        GetCurrentTimeServiceUseCase: GetCurrentTimeUseCase,  // ✅ Inyección
    }
}
```

## Reglas de Importación

### ✅ Reglas Correctas

1. **Domain importa solo estándar de Go**
2. **Application importa solo Domain**
3. **Infrastructure importa solo Domain** (para implementar puertos)
4. **Presentation importa Domain y Application**
5. **Container resuelve todas las dependencias**

### ❌ Violaciones Comunes

1. **Domain importando Infrastructure**:
```go
// ❌ ERROR
import "genesis/pos/infrastructure/db/postgres"
```

2. **Application importando Infrastructure directamente**:
```go
// ❌ ERROR
import "genesis/pos/infrastructure/repositories_impl"
```

3. **Infrastructure importando Application**:
```go
// ❌ ERROR
import "genesis/pos/application/use_case"
```

4. **Domain importando Presentation**:
```go
// ❌ ERROR
import "genesis/pos/presentation/dtos"
```

## Verificación de Dependencias

Para verificar que las dependencias son correctas, puedes usar:

```bash
# Verificar imports de un paquete
go list -f '{{.Imports}}' ./domain/...

# Verificar que domain no importa infrastructure
go list -f '{{.Imports}}' ./domain/... | grep infrastructure
# No debe retornar nada
```

## Inversión de Dependencias

La arquitectura hexagonal usa **Inversión de Dependencias** (Dependency Inversion Principle):

### Sin Inversión (❌ Incorrecto):
```
Application → Infrastructure (dependencia directa)
```

### Con Inversión (✅ Correcto):
```
Application → Domain Port (interfaz)
Infrastructure → Domain Port (implementa interfaz)
```

**Resultado**: Application no depende de Infrastructure, ambos dependen de abstracciones (puertos) en Domain.

## Ejemplo Completo de Flujo

```
Request HTTP
    ↓
[Presentation] Handler recibe request
    ↓ (usa)
[Application] Service orquesta caso de uso
    ↓ (usa puerto)
[Domain] Puerto (interfaz) GetTimeRepository
    ↑ (implementa)
[Infrastructure] Repository Implementation
    ↓ (accede)
Base de Datos
```

**Código**:
```go
// Presentation
func GetTime(ctx *gin.Context) {
    service := container.GetTimeService  // Obtiene de container
    response := service.ExecuteService()
}

// Application
func (s *TimeServiceImpl) ExecuteService() {
    time, err := s.useCase.ExecuteUseCase()  // Usa caso de uso
}

// Application
func (uc *TimeUseCaseImpl) ExecuteUseCase() {
    time, err := uc.repo.ExecuteImpl()  // Usa puerto (interfaz)
}

// Domain (Puerto)
type GetTimeRepository interface {
    ExecuteImpl() (time.Time, error)
}

// Infrastructure (Implementación)
func (r *GetTimeWithPostgresImpl) ExecuteImpl() (time.Time, error) {
    // Implementación con PostgreSQL
}
```

## Resumen

| Capa | Puede Importar | NO Puede Importar |
|------|---------------|-------------------|
| **Domain** | Solo stdlib | Todo lo demás |
| **Application** | Domain, stdlib | Infrastructure, Presentation |
| **Infrastructure** | Domain, frameworks | Application, Presentation |
| **Presentation** | Domain, Application, frameworks | Infrastructure (directamente) |

**Regla de Oro**: Si una capa interna importa una capa externa, **estás violando la arquitectura**.

