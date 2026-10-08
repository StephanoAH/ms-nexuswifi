# Decisiones Arquitectónicas Iniciales (ADR)

## Introducción

Este documento registra las decisiones arquitectónicas iniciales del proyecto. Las ADRs (Architecture Decision Records) documentan el contexto, la decisión y las consecuencias de elecciones arquitectónicas importantes.

## ADR-001: Arquitectura Hexagonal

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos una arquitectura que permitiera:
- Separar la lógica de negocio de los detalles técnicos
- Facilitar el cambio de tecnologías (bases de datos, frameworks)
- Hacer el código testeable y mantenible
- Escalar la aplicación sin acoplamiento

**Decisión**:
Implementar **Arquitectura Hexagonal** (Ports and Adapters) con cuatro capas:
- `domain/`: Núcleo con entidades, puertos y value objects
- `application/`: Casos de uso y servicios
- `infrastructure/`: Implementaciones técnicas (DB, APIs externas)
- `presentation/`: Adaptadores de entrada (HTTP, handlers)

**Consecuencias**:
- ✅ **Positivas**:
  - Código desacoplado y mantenible
  - Fácil cambiar implementaciones (ej: PostgreSQL → MongoDB)
  - Testing simplificado (mockear puertos)
  - Escalabilidad mejorada
- ⚠️ **Negativas**:
  - Más archivos y estructura compleja
  - Curva de aprendizaje inicial
  - Overhead en proyectos pequeños

## ADR-002: Inyección de Dependencias Manual

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos un mecanismo de inyección de dependencias. Opciones:
1. Container manual (patrón actual)
2. Librería de DI (wire, dig, fx)
3. Constructor injection puro

**Decisión**:
Usar **Container Pattern manual** en `presentation/container/`:
- Variables globales para instancias
- Funciones `ResolveXxx()` con patrón Singleton
- Inicialización centralizada en `StartContainer()`

**Consecuencias**:
- ✅ **Positivas**:
  - Sin dependencias externas
  - Control total sobre resolución
  - Simple y directo
  - Fácil de entender
- ⚠️ **Negativas**:
  - Variables globales (riesgo de estado compartido)
  - No hay validación de dependencias en tiempo de compilación
  - Mantenimiento manual

**Alternativas Consideradas**:
- **wire**: Generación de código, más seguro pero requiere build step
- **dig**: DI basada en reflection, más flexible pero menos type-safe
- **fx**: Framework completo, más overhead para este proyecto

## ADR-003: Naming Conventions

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos convenciones claras de naming para mantener consistencia.

**Decisión**:
- **Interfaces**: Sin sufijo (ej: `GetUserUseCase`)
- **Implementaciones**: Sufijo `Impl` (ej: `GetUserUseCaseImpl`)
- **Repositorios con DB**: `With[DB]Impl` (ej: `GetUserWithPostgresImpl`)
- **Métodos**: `ExecuteUseCase()`, `ExecuteImpl()`, `ExecuteService()`

**Consecuencias**:
- ✅ **Positivas**:
  - Consistencia en todo el código
  - Fácil identificar tipos (interfaz vs implementación)
  - Naming descriptivo y claro
- ⚠️ **Negativas**:
  - Nombres más largos
  - Requiere disciplina del equipo

## ADR-004: Validación con ozzo-validation

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos validar requests HTTP. Opciones:
1. Validación manual en handlers
2. Librería de validación (ozzo-validation, validator)
3. Schemas con métodos Validate()

**Decisión**:
Usar **ozzo-validation** en schemas (`presentation/schema/`):
- Structs con método `Validate()` que implementa `Validatable`
- Middlewares genéricos que usan la interfaz
- Validación declarativa y reutilizable

**Consecuencias**:
- ✅ **Positivas**:
  - Validación declarativa y clara
  - Reutilizable en múltiples endpoints
  - Middlewares genéricos
  - Fácil agregar nuevas validaciones
- ⚠️ **Negativas**:
  - Dependencia externa
  - Curva de aprendizaje de la librería

**Alternativas Consideradas**:
- **validator**: Más popular pero menos flexible
- **Validación manual**: Más control pero más código repetitivo

## ADR-005: Framework HTTP - Gin

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos un framework HTTP para la capa de presentación.

**Decisión**:
Usar **Gin** como framework HTTP:
- Framework popular y estable
- Performance buena
- Middleware support
- Binding automático de requests

**Consecuencias**:
- ✅ **Positivas**:
  - Framework maduro y estable
  - Buena documentación
  - Comunidad activa
  - Performance excelente
- ⚠️ **Negativas**:
  - Acoplamiento a Gin en capa de presentación
  - Cambiar framework requiere refactor de presentación

**Nota**: El acoplamiento está limitado a `presentation/`, no afecta otras capas gracias a la arquitectura hexagonal.

## ADR-006: Bases de Datos - PostgreSQL y MongoDB

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
El proyecto necesita soportar múltiples bases de datos.

**Decisión**:
Soporte para **PostgreSQL** y **MongoDB**:
- PostgreSQL con `pgx/v5`
- MongoDB con `mongo-driver/v2`
- Drivers encapsulados en `infrastructure/db/`
- Repositorios específicos por DB (`WithPostgres`, `WithMongo`)

**Consecuencias**:
- ✅ **Positivas**:
  - Flexibilidad para elegir DB según necesidad
  - Fácil agregar nuevas bases de datos
  - Repositorios desacoplados de drivers
- ⚠️ **Negativas**:
  - Más código (implementaciones por DB)
  - Mantenimiento de múltiples drivers

## ADR-007: Manejo de Errores - Propagación Simple

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos una estrategia clara de manejo de errores.

**Decisión**:
**Propagación simple** de errores:
- Retornar `error` como último valor
- Propagación directa sin modificación en la mayoría de casos
- Logging opcional antes de retornar
- Transformación a HTTP solo en handlers

**Consecuencias**:
- ✅ **Positivas**:
  - Simple y directo
  - No pierde información del error original
  - Fácil de seguir el flujo
- ⚠️ **Negativas**:
  - No agrega contexto al error
  - Puede ser difícil rastrear origen del error

**Mejoras Futuras Consideradas**:
- Errores tipados (`var ErrNotFound = errors.New(...)`)
- Contexto con `fmt.Errorf` y `%w`
- Logging estructurado

## ADR-008: Testing - TDD sin Librerías Externas Inicialmente

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
El proyecto actualmente no tiene tests. Necesitábamos definir estrategia.

**Decisión**:
**TDD con testing estándar** inicialmente:
- Usar `testing` package estándar
- Mocks manuales
- Considerar `testify` y `testcontainers` para el futuro

**Consecuencias**:
- ✅ **Positivas**:
  - Sin dependencias externas
  - Control total sobre mocks
  - Fácil de entender
- ⚠️ **Negativas**:
  - Más código para mocks
  - Menos assertions helpers
  - Tests de integración más complejos

**Futuro**:
- Evaluar `testify` para assertions
- Evaluar `testcontainers` para tests de integración

## ADR-009: Estructura de Respuestas HTTP

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos estructuras consistentes para respuestas HTTP.

**Decisión**:
Dos estructuras principales:
- `ResponseHttp`: Para respuestas simples (`code`, `success`, `message`)
- `ResponseJson`: Para respuestas estructuradas (`status`, `message`, `data`)

**Consecuencias**:
- ✅ **Positivas**:
  - Consistencia en respuestas
  - Fácil de usar en handlers
  - Estructura clara para clientes
- ⚠️ **Negativas**:
  - Dos estructuras diferentes (puede ser confuso)
  - No hay validación de estructura en tiempo de compilación

## ADR-010: Logging - log Estándar

**Fecha**: Inicial del proyecto

**Estado**: Aceptada

**Contexto**:
Necesitábamos logging para debugging y monitoreo.

**Decisión**:
Usar `log` package estándar de Go:
- `log.Println()` para errores
- Configuración básica en `main.go`

**Consecuencias**:
- ✅ **Positivas**:
  - Sin dependencias externas
  - Simple y directo
  - Suficiente para desarrollo
- ⚠️ **Negativas**:
  - No hay niveles de log (INFO, WARN, ERROR)
  - No hay logging estructurado
  - Limitado para producción

**Futuro**:
- Evaluar `zerolog` o `zap` para logging estructurado
- Agregar niveles de log
- Integración con sistemas de monitoreo

## Resumen de ADRs

| ADR | Decisión | Estado |
|-----|----------|--------|
| ADR-001 | Arquitectura Hexagonal | ✅ Aceptada |
| ADR-002 | DI Manual (Container) | ✅ Aceptada |
| ADR-003 | Naming Conventions | ✅ Aceptada |
| ADR-004 | Validación con ozzo-validation | ✅ Aceptada |
| ADR-005 | Framework HTTP - Gin | ✅ Aceptada |
| ADR-006 | Bases de Datos - PostgreSQL y MongoDB | ✅ Aceptada |
| ADR-007 | Manejo de Errores - Propagación Simple | ✅ Aceptada |
| ADR-008 | Testing - TDD Estándar | ✅ Aceptada |
| ADR-009 | Estructura de Respuestas HTTP | ✅ Aceptada |
| ADR-010 | Logging - log Estándar | ✅ Aceptada |

## Proceso para Nuevas ADRs

Cuando se necesite tomar una decisión arquitectónica importante:

1. **Crear nuevo archivo ADR**: `adr-[número].md`
2. **Documentar**:
   - Contexto
   - Decisión
   - Consecuencias
   - Alternativas consideradas
3. **Revisar con el equipo**
4. **Actualizar este documento** con referencia

## Notas

- Estas ADRs reflejan el estado actual del proyecto
- Pueden evolucionar según necesidades
- Nuevas decisiones deben documentarse como ADRs
- Cambios importantes requieren revisión de ADRs afectadas

