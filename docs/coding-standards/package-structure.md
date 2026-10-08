# Estructura de Paquetes

## Introducción

Este documento describe la organización de paquetes del proyecto y las convenciones para mantener la estructura consistente.

## Estructura General

```
genesis/pos/
├── domain/                    # Capa de Dominio
│   ├── entities/              # Entidades del negocio
│   ├── ports/                 # Puertos (interfaces)
│   │   ├── application/       # Puertos para aplicación
│   │   │   ├── use_case/      # Puertos de casos de uso
│   │   │   └── service_application/ # Puertos de servicios
│   │   └── infrastructure/    # Puertos para infraestructura
│   └── value_object/          # Value objects y constantes
│       └── constants/         # Constantes del dominio
├── application/               # Capa de Aplicación
│   ├── use_case/              # Casos de uso
│   └── service_application/   # Servicios de aplicación
├── infrastructure/            # Capa de Infraestructura
│   ├── repositories_impl/     # Implementaciones de repositorios
│   ├── db/                    # Drivers de base de datos
│   │   ├── postgres/          # PostgreSQL
│   │   │   └── drivers/       # Drivers de conexión
│   │   └── mongo/             # MongoDB
│   │       └── drivers/       # Drivers de conexión
│   └── external_http/         # Clientes HTTP externos
│       └── net_http/          # Implementación con net/http
└── presentation/              # Capa de Presentación
    ├── handlers/              # Handlers HTTP
    ├── routers/               # Configuración de rutas
    ├── middleware/            # Middlewares
    ├── dtos/                  # Data Transfer Objects
    ├── schema/                 # Schemas de validación
    ├── container/              # Container de inyección de dependencias
    └── boostrap/              # Bootstrap del servidor
```

## Domain Layer

### `domain/entities/`

**Propósito**: Contiene las entidades del dominio de negocio.

**Contenido**:
- Estructuras que representan conceptos del negocio
- Tags JSON y DB para serialización
- Lógica de negocio asociada (cuando aplica)

**Ejemplos**:
- `response_current_time.go`: `ResponseGetCurrentTime`
- `basic_auth.go`: `BasicAuth`
- `sales_mongo.go`: `VentaDetallada`
- `format_response.go`: `ResponseHttp`, `ResponseJson`

**Reglas**:
- ✅ Solo entidades del dominio
- ✅ Sin dependencias de otras capas
- ✅ Nombres descriptivos del negocio

### `domain/ports/`

**Propósito**: Define los contratos (interfaces) que otras capas deben implementar.

#### `domain/ports/application/`

**Propósito**: Puertos para la capa de aplicación.

**Subdirectorios**:
- `use_case/`: Interfaces de casos de uso
- `service_application/`: Interfaces de servicios

**Ejemplos**:
- `use_case/get_current_time_use_case.go`: `GetCurrentTimeUseCase`
- `service_application/get_current_time_service.go`: `GetCurrentTimeService`

#### `domain/ports/infrastructure/`

**Propósito**: Puertos para la capa de infraestructura.

**Contenido**:
- Interfaces de repositorios
- Interfaces de conexiones de base de datos
- Interfaces de servicios externos

**Ejemplos**:
- `get_current_time_repository.go`: `GetCurrentTimeRepository`
- `database_connection_repository.go`: `DatabaseConnectionPgxInterface`
- `basic_auth_repository.go`: `BasicAuthRepositoryPort`

### `domain/value_object/`

**Propósito**: Objetos de valor y constantes del dominio.

**Estructura**:
- `constants/`: Constantes del dominio

**Ejemplos**:
- `constants/running_server.go`: Constantes del servidor
- `constants/url_service.go`: URLs de servicios

## Application Layer

### `application/use_case/`

**Propósito**: Implementa los casos de uso del negocio.

**Contenido**:
- Implementaciones de casos de uso
- Cada archivo corresponde a un caso de uso específico

**Patrón de Naming**:
- `get_[entity]_use_case.go`: Implementación del caso de uso

**Ejemplos**:
- `get_current_time_use_case.go`: `GetCurrentTimeUseCaseImpl`
- `get_sales_use_case.go`: `GetSalesUseCaseImpl`
- `get_credentials_use_case.go`: `GetCredentialsUseCaseImpl`

**Reglas**:
- ✅ Implementa puertos de `domain/ports/application/use_case/`
- ✅ Solo depende de `domain/`
- ✅ Un caso de uso por archivo

### `application/service_application/`

**Propósito**: Servicios que orquestan múltiples casos de uso.

**Contenido**:
- Implementaciones de servicios de aplicación
- Coordinan operaciones complejas

**Ejemplos**:
- `get_current_time_service.go`: `GetCurrentTimeServiceImpl`

**Reglas**:
- ✅ Implementa puertos de `domain/ports/application/service_application/`
- ✅ Puede llamar a múltiples casos de uso
- ✅ Transforma respuestas de casos de uso

## Infrastructure Layer

### `infrastructure/repositories_impl/`

**Propósito**: Implementaciones concretas de repositorios.

**Contenido**:
- Implementaciones de puertos de `domain/ports/infrastructure/`
- Lógica de acceso a datos

**Patrón de Naming**:
- `get_[entity]_with_postgres.go`: Implementación con PostgreSQL
- `get_[entity]_with_mongo.go`: Implementación con MongoDB
- `[entity]_repository.go`: Implementación genérica

**Ejemplos**:
- `get_current_time_with_postgres.go`: `GetCurrentTimeWithPostgresImpl`
- `get_sales_with_mongo.go`: `GetSalesWithMongoImpl`
- `basic_auth_repository.go`: `BasicAuthRepositoryImpl`
- `get_environment.go`: `GetEnvironmentImpl`

**Reglas**:
- ✅ Implementa puertos de `domain/ports/infrastructure/`
- ✅ Un repositorio por archivo
- ✅ Naming indica la tecnología usada (`WithPostgres`, `WithMongo`)

### `infrastructure/db/`

**Propósito**: Drivers y configuraciones de bases de datos.

#### `infrastructure/db/postgres/`

**Estructura**:
- `drivers/`: Implementaciones de drivers
  - `config_connection_pgx.go`: Configuración de conexión
- `config_database_connection_pgx_driver.go`: Driver de conexión

**Ejemplos**:
- `drivers/config_connection_pgx.go`: `ConfigConnectionPgx`
- `config_database_connection_pgx_driver.go`: `ConfigDatabaseConnectionPgxDriver`

#### `infrastructure/db/mongo/`

**Estructura**:
- `drivers/`: Implementaciones de drivers
  - `config_connection_mongo.go`: Configuración de conexión
- `config_database_connection_mongo_driver.go`: Driver de conexión

**Ejemplos**:
- `drivers/config_connection_mongo.go`: `ConfigConnectionMongo`
- `config_database_connection_mongo_driver.go`: `ConfigDatabaseConnectionMongoDriver`

### `infrastructure/external_http/`

**Propósito**: Cliente HTTP para APIs externas.

**Estructura**:
- `http_client_interface.go`: Interfaz del cliente
- `http_request.go`: Estructura de request (si existe)
- `http_response.go`: Estructura de response (si existe)
- `net_http/`: Implementación con `net/http`
  - `net_http_client.go`: `NetHTTPClient`

**Reglas**:
- ✅ Interfaz en el nivel superior
- ✅ Implementaciones en subdirectorios por tecnología

## Presentation Layer

### `presentation/handlers/`

**Propósito**: Handlers HTTP que procesan requests.

**Contenido**:
- Funciones que reciben `*gin.Context`
- Transforman requests en llamadas a servicios

**Patrón de Naming**:
- `get_[entity].go`: Handler para obtener entidad
- `[action]_[entity].go`: Handler para acción específica

**Ejemplos**:
- `get_current_time.go`: `GetCurrentTime`

**Reglas**:
- ✅ Una función handler por archivo
- ✅ Nombres descriptivos de la acción

### `presentation/routers/`

**Propósito**: Configuración de rutas HTTP.

**Contenido**:
- Funciones que configuran grupos de rutas
- Aplican middlewares a rutas

**Patrón de Naming**:
- `[resource].go`: Configuración de rutas para un recurso

**Ejemplos**:
- `configuration.go`: `CurrentTime` (configura rutas de current-time)

**Reglas**:
- ✅ Una función por grupo de rutas
- ✅ Aplica middlewares de validación

### `presentation/middleware/`

**Propósito**: Middlewares para cross-cutting concerns.

**Contenido**:
- Middlewares de autenticación
- Middlewares de validación
- Middlewares de logging
- Middlewares de recovery

**Ejemplos**:
- `jwt_middleware.go`: `JWTMiddleware`
- `basic_auth_middleware.go`: `BasicAuthMiddleware`
- `validate_routes.go`: `ValidateBody`, `ValidateQuery`, `ValidateParams`
- `logger_middleware.go`: `LoggerMiddleware`
- `recovery_middleware.go`: `RecoveryMiddleware`

**Reglas**:
- ✅ Un middleware por archivo
- ✅ Retornan `gin.HandlerFunc`

### `presentation/dtos/`

**Propósito**: Data Transfer Objects para requests HTTP.

**Contenido**:
- Estructuras para recibir datos de requests
- Tags JSON para binding

**Ejemplos**:
- `movement_dto.go`: `MovementDto`

**Reglas**:
- ✅ Solo para datos de entrada/salida HTTP
- ✅ No contienen lógica de negocio

### `presentation/schema/`

**Propósito**: Schemas de validación para requests.

**Contenido**:
- Estructuras con métodos `Validate()`
- Usan `ozzo-validation`

**Ejemplos**:
- `current_time.go`: `CurrentTimeBody`, `CurrentTimeParams`, `CurrentTimeQuery`

**Reglas**:
- ✅ Implementan interfaz `Validatable`
- ✅ Método `Validate()` con `ozzo-validation`

### `presentation/container/`

**Propósito**: Container de inyección de dependencias.

**Contenido**:
- Variables globales para instancias
- Funciones `ResolveXxx()` para resolver dependencias
- Función `StartContainer()` para inicialización

**Ejemplos**:
- `container.go`: Todas las resoluciones de dependencias

**Reglas**:
- ✅ Patrón Singleton para instancias
- ✅ Inicialización centralizada

### `presentation/boostrap/`

**Propósito**: Bootstrap y configuración del servidor.

**Contenido**:
- Configuración del servidor HTTP
- Carga de middlewares globales
- Configuración de rutas

**Ejemplos**:
- `server.go`: `RunServer()`, `loadEngine()`
- `upload_env.go`: `UploadEnv()`

**Nota**: `boostrap` tiene un typo (debería ser `bootstrap`), pero se mantiene por consistencia.

## Reglas Generales

### ✅ DEBE

1. **Un concepto por paquete**: Cada paquete tiene un propósito claro
2. **Naming consistente**: Seguir patrones establecidos
3. **Separación de concerns**: Cada capa en su directorio
4. **Dependencias hacia adentro**: Capas externas dependen de internas

### ❌ NO DEBE

1. **Mezclar responsabilidades**: No poner handlers en `use_case/`
2. **Dependencias incorrectas**: No importar capas externas desde internas
3. **Naming inconsistente**: No mezclar convenciones
4. **Archivos fuera de lugar**: No poner entidades en `infrastructure/`

## Agregar Nuevos Paquetes

### Nuevo Caso de Uso

1. Crear puerto en `domain/ports/application/use_case/`
2. Crear implementación en `application/use_case/`
3. Agregar resolución en `presentation/container/`

### Nuevo Repositorio

1. Crear puerto en `domain/ports/infrastructure/`
2. Crear implementación en `infrastructure/repositories_impl/`
3. Agregar resolución en `presentation/container/`

### Nuevo Handler

1. Crear handler en `presentation/handlers/`
2. Crear rutas en `presentation/routers/`
3. Agregar al bootstrap si es necesario

## Resumen

La estructura de paquetes sigue la arquitectura hexagonal:

- **Domain**: Entidades, puertos, value objects
- **Application**: Casos de uso, servicios
- **Infrastructure**: Repositorios, drivers, clientes externos
- **Presentation**: Handlers, routers, middlewares, DTOs, container, bootstrap

Cada paquete tiene un propósito claro y sigue convenciones de naming consistentes.

