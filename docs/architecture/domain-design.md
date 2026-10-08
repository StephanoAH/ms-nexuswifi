# Diseño del Dominio

## Introducción

El dominio es el **núcleo** de la arquitectura hexagonal. Contiene toda la lógica de negocio pura, sin dependencias de frameworks, bases de datos o sistemas externos.

## Estructura del Dominio

```
domain/
├── entities/          # Entidades del dominio
├── ports/             # Puertos (interfaces)
│   ├── application/   # Puertos para capa de aplicación
│   └── infrastructure/ # Puertos para capa de infraestructura
└── value_object/      # Objetos de valor y constantes
```

## Entidades (Entities)

### Concepto

Las entidades representan **conceptos del negocio** que tienen identidad y ciclo de vida. Son estructuras de datos que encapsulan lógica de dominio.

### Ubicación

`domain/entities/`

### Características

1. **Inmutabilidad cuando es posible**: Preferir estructuras que no cambien
2. **Lógica de negocio**: Pueden tener métodos que implementen reglas de negocio
3. **Sin dependencias externas**: No importan paquetes de infraestructura o presentación
4. **Serialización**: Usan tags JSON para serialización

### Ejemplos en el Proyecto

#### ResponseGetCurrentTime

```go
type ResponseGetCurrentTime struct {
    StatusCode int       `json:"StatusCode"`
    Message    string    `json:"Message"`
    Result     time.Time `json:"CurrentTime"`
}
```

**Uso**: Representa la respuesta estructurada de un caso de uso que retorna la hora actual.

#### BasicAuth

```go
type BasicAuth struct {
    Username string `json:"username" db:"identificador_conexion"`
    Password string `json:"password" db:"clave_conexion"`
}
```

**Uso**: Representa credenciales de autenticación básica. Nota los tags `db:` para mapeo a base de datos.

#### VentaDetallada

```go
type VentaDetallada struct {
    ID struct {
        Oid string `json:"$oid"`
    } `json:"_id"`
    IDMovimiento       int      `json:"id_movimiento"`
    // ... más campos
}
```

**Uso**: Representa una venta detallada desde MongoDB. Incluye mapeo específico de MongoDB (`$oid`).

#### ResponseHttp

```go
type ResponseHttp struct {
    Code    int    `json:"code"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

**Uso**: Respuesta HTTP genérica para handlers.

#### ResponseJson

```go
type ResponseJson struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

**Uso**: Respuesta JSON estructurada con datos opcionales.

### Convenciones para Entidades

1. **Naming**: PascalCase, nombres descriptivos del dominio
2. **Tags JSON**: Siempre incluir tags `json:` para serialización
3. **Tags DB**: Incluir tags `db:` cuando se mapean directamente a tablas
4. **Documentación**: Agregar comentarios cuando el propósito no es obvio

## Puertos (Ports)

### Concepto

Los puertos son **interfaces** que definen contratos. Representan "qué necesita" el dominio, no "cómo se implementa".

### Ubicación

`domain/ports/`

### Tipos de Puertos

#### 1. Application Ports

**Ubicación**: `domain/ports/application/`

Definen contratos para la capa de aplicación:

- **Use Case Ports**: `domain/ports/application/use_case/`
- **Service Ports**: `domain/ports/application/service_application/`

**Ejemplo - Use Case Port**:
```go
package use_case

import "time"

type GetCurrentTimeUseCase interface {
    ExecuteUseCase() (time.Time, error)
}
```

**Ejemplo - Service Port**:
```go
package service_application

import "genesis/pos/domain/entities"

type GetCurrentTimeService interface {
    ExecuteService() entities.ResponseGetCurrentTime
}
```

#### 2. Infrastructure Ports

**Ubicación**: `domain/ports/infrastructure/`

Definen contratos para la capa de infraestructura:

- **Repository Ports**: Contratos para acceso a datos
- **Connection Ports**: Contratos para conexiones de base de datos
- **External Service Ports**: Contratos para servicios externos

**Ejemplo - Repository Port**:
```go
package infrastructure

import "time"

type GetCurrentTimeRepository interface {
    ExecuteImpl() (time.Time, error)
}
```

**Ejemplo - Connection Port**:
```go
package infrastructure

import "genesis/pos/infrastructure/db/postgres"

type DatabaseConnectionPgxInterface interface {
    GetDatabaseConnectionPgx() *postgres.ConfigDatabaseConnectionPgxDriver
}
```

**Ejemplo - External Service Port**:
```go
package infrastructure

type GetEnvironmentRepository interface {
    GetEnvironmentString(key string) string
    GetEnvironmentInt(key string) (int, error)
    GetEnvironmentBool(key string) (bool, error)
}
```

### Convenciones para Puertos

1. **Naming de Interfaces**:
   - **Sin sufijo**: `GetCurrentTimeUseCase` (no `GetCurrentTimeUseCaseInterface`)
   - **Con sufijo "Port" cuando hay ambigüedad**: `BasicAuthRepositoryPort`
   - **Con sufijo "Interface" para conexiones**: `DatabaseConnectionPgxInterface`

2. **Métodos**:
   - **Use Cases**: `ExecuteUseCase()`
   - **Repositories**: `ExecuteImpl()` o métodos específicos como `GetCredentials()`
   - **Services**: `ExecuteService()`

3. **Retornos**:
   - Siempre incluir `error` como último retorno
   - Retornar entidades del dominio, no DTOs

4. **Ubicación de Interfaces**:
   - Interfaces de aplicación → `domain/ports/application/`
   - Interfaces de infraestructura → `domain/ports/infrastructure/`

## Value Objects

### Concepto

Los value objects son objetos **inmutables** que representan conceptos del dominio sin identidad propia. Se comparan por valor, no por referencia.

### Ubicación

`domain/value_object/`

### Estructura Actual

```
value_object/
└── constants/
    ├── running_server.go
    └── url_service.go
```

### Ejemplo

```go
package constants

var ID_POS int64 = 0

const DB_CON string = "postgresql://..."
const HOST_PORT string = ":3001"
const currentEnvironment = "ws.desarrollo.genexis.com.co"
```

### Convenciones para Value Objects

1. **Inmutabilidad**: Los value objects no deben cambiar después de crearse
2. **Comparación por valor**: Dos value objects con el mismo valor son iguales
3. **Sin identidad**: No tienen ID único
4. **Validación**: Deben validarse en construcción

## Principios de Diseño del Dominio

### 1. Rich Domain Model

Preferir entidades con lógica de negocio sobre anemic domain models:

```go
// ✅ Bueno: Lógica en la entidad
func (u *User) IsActive() bool {
    return u.Status == "active" && !u.DeletedAt.IsZero()
}

// ❌ Evitar: Entidad anémica sin lógica
type User struct {
    Status string
    DeletedAt time.Time
}
// Lógica en otro lugar
```

### 2. Ubiquitous Language

Usar el lenguaje del dominio en nombres de entidades, métodos y propiedades:

- ✅ `VentaDetallada` (lenguaje del negocio)
- ❌ `SalesDetailRecord` (lenguaje técnico)

### 3. Encapsulación

Encapsular invariantes y reglas de negocio dentro de las entidades:

```go
type Order struct {
    Items []Item
    Total float64
}

func (o *Order) AddItem(item Item) error {
    if item.Price <= 0 {
        return errors.New("item price must be positive")
    }
    o.Items = append(o.Items, item)
    o.Total += item.Price
    return nil
}
```

### 4. Separación de Concerns

- **Entidades**: Lógica de negocio y datos
- **Puertos**: Contratos e interfaces
- **Value Objects**: Conceptos inmutables

## Reglas del Dominio

### ✅ DEBE

1. El dominio **NO debe** importar paquetes de `infrastructure/` o `presentation/`
2. Las entidades deben ser estructuras simples con tags de serialización
3. Los puertos deben definir solo lo necesario, no detalles de implementación
4. Usar tipos del dominio, no tipos primitivos cuando sea apropiado

### ❌ NO DEBE

1. El dominio **NO debe** conocer frameworks (Gin, pgx, mongo-driver, etc.)
2. Las entidades **NO deben** tener dependencias de infraestructura
3. Los puertos **NO deben** exponer detalles de implementación
4. **NO usar** DTOs en el dominio (solo entidades)

## Ejemplo Completo: Flujo de Diseño

### 1. Identificar Concepto del Dominio

**Necesidad**: Obtener credenciales de autenticación

### 2. Crear Entidad

```go
// domain/entities/basic_auth.go
type BasicAuth struct {
    Username string `json:"username" db:"identificador_conexion"`
    Password string `json:"password" db:"clave_conexion"`
}
```

### 3. Definir Puerto

```go
// domain/ports/infrastructure/basic_auth_repository.go
type BasicAuthRepositoryPort interface {
    GetCredentials(integrationID int64, username string) (entities.BasicAuth, error)
}
```

### 4. Implementar en Infraestructura

```go
// infrastructure/repositories_impl/basic_auth_repository.go
type BasicAuthRepositoryImpl struct {
    DB infrastructure.DatabaseConnectionPgxInterface
}

func (b *BasicAuthRepositoryImpl) GetCredentials(...) (entities.BasicAuth, error) {
    // Implementación con PostgreSQL
}
```

## Siguientes Pasos

- Ver [Capas y Dependencias](./layers-and-dependencies.md)
- Ver [Estándares de Código](../coding-standards/go-standards.md)
- Ver [Convenciones de Nombres](../coding-standards/naming-conventions.md)

