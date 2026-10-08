# Visión General de la Arquitectura

## Introducción

Este proyecto implementa una **Arquitectura Hexagonal** (también conocida como Arquitectura de Puertos y Adaptadores) en Go. Esta arquitectura separa la lógica de negocio de los detalles técnicos, permitiendo que la aplicación sea independiente de frameworks, bases de datos y sistemas externos.

## Principios Fundamentales

1. **Separación de Responsabilidades**: Cada capa tiene una responsabilidad clara y bien definida
2. **Inversión de Dependencias**: Las capas internas no dependen de las externas
3. **Independencia Tecnológica**: El dominio no conoce detalles de implementación
4. **Testabilidad**: Cada componente puede ser probado de forma aislada

## Estructura de Capas

El proyecto está organizado en cuatro capas principales:

```
genesis/pos/
├── domain/              # Capa de Dominio (Núcleo)
├── application/         # Capa de Aplicación
├── infrastructure/      # Capa de Infraestructura (Adaptadores de Salida)
└── presentation/        # Capa de Presentación (Adaptadores de Entrada)
```

### 1. Domain (Dominio)

**Ubicación**: `domain/`

**Responsabilidades**:
- Contiene las entidades del negocio
- Define los puertos (interfaces) que representan contratos
- Incluye value objects y constantes
- **NO depende de ninguna otra capa**

**Estructura**:
- `entities/`: Entidades del dominio
- `ports/`: Interfaces que definen contratos
  - `application/`: Puertos para casos de uso y servicios
  - `infrastructure/`: Puertos para repositorios y conexiones
- `value_object/`: Objetos de valor y constantes

**Ejemplo de Entidad**:
```go
type ResponseGetCurrentTime struct {
    StatusCode int       `json:"StatusCode"`
    Message    string    `json:"Message"`
    Result     time.Time `json:"CurrentTime"`
}
```

### 2. Application (Aplicación)

**Ubicación**: `application/`

**Responsabilidades**:
- Implementa los casos de uso del negocio
- Orquesta la lógica de aplicación
- Coordina entre repositorios y servicios
- **Depende solo de `domain/`**

**Estructura**:
- `use_case/`: Casos de uso que implementan la lógica de negocio
- `service_application/`: Servicios que orquestan múltiples casos de uso

**Flujo típico**:
1. Handler llama a Service
2. Service llama a Use Case
3. Use Case llama a Repository (a través de puerto)

### 3. Infrastructure (Infraestructura)

**Ubicación**: `infrastructure/`

**Responsabilidades**:
- Implementa los adaptadores de salida
- Proporciona implementaciones concretas de los puertos
- Maneja conexiones a bases de datos, APIs externas, etc.
- **Depende de `domain/` y `application/`**

**Estructura**:
- `repositories_impl/`: Implementaciones de repositorios
- `db/`: Drivers y configuraciones de bases de datos
  - `postgres/`: Configuración PostgreSQL
  - `mongo/`: Configuración MongoDB
- `external_http/`: Clientes HTTP para APIs externas

### 4. Presentation (Presentación)

**Ubicación**: `presentation/`

**Responsabilidades**:
- Implementa los adaptadores de entrada
- Maneja HTTP, routing, middlewares
- Transforma requests en llamadas a la capa de aplicación
- **Depende de `domain/` y `application/`**

**Estructura**:
- `handlers/`: Manejadores HTTP (endpoints)
- `routers/`: Configuración de rutas
- `middleware/`: Middlewares (JWT, validación, logging, etc.)
- `dtos/`: Data Transfer Objects para requests/responses
- `schema/`: Schemas de validación
- `container/`: Inyección de dependencias
- `boostrap/`: Configuración e inicio del servidor

## Flujo de Datos

```
Request HTTP
    ↓
[Presentation Layer]
    ├── Middleware (validación, autenticación, logging)
    ├── Handler (recibe request, valida DTO)
    └── Container (resuelve dependencias)
    ↓
[Application Layer]
    ├── Service (orquesta casos de uso)
    └── Use Case (lógica de negocio)
    ↓
[Domain Layer]
    └── Port (interfaz del repositorio)
    ↓
[Infrastructure Layer]
    └── Repository Implementation (acceso a datos)
    ↓
Base de Datos / API Externa
```

## Ventajas de esta Arquitectura

1. **Mantenibilidad**: Código organizado y fácil de entender
2. **Testabilidad**: Cada capa puede probarse independientemente
3. **Flexibilidad**: Fácil cambiar implementaciones (ej: cambiar de PostgreSQL a MongoDB)
4. **Escalabilidad**: Fácil agregar nuevas funcionalidades sin afectar existentes
5. **Independencia**: El dominio no cambia si cambiamos frameworks o tecnologías

## Reglas de Dependencias

- ✅ `domain` → **NO depende de nadie**
- ✅ `application` → Solo depende de `domain`
- ✅ `infrastructure` → Depende de `domain` (implementa puertos)
- ✅ `presentation` → Depende de `domain` y `application`
- ❌ **NUNCA**: `domain` o `application` deben importar `infrastructure` o `presentation`

## Punto de Entrada

El archivo `main.go` es el punto de entrada de la aplicación:

```go
func main() {
    boostrap.UploadEnv()           // Carga variables de entorno
    container.StartContainer()     // Inicializa dependencias
    boostrap.RunServer()           // Inicia el servidor HTTP
}
```

## Siguientes Pasos

- Ver [Estructura Hexagonal Detallada](./hexagonal-structure.md)
- Ver [Diseño del Dominio](./domain-design.md)
- Ver [Capas y Dependencias](./layers-and-dependencies.md)

