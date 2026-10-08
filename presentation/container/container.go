package container

// StartContainer inicializa todas las dependencias del contenedor de inyección de dependencias.
// Orquesta la inicialización en el siguiente orden:
// 1. Conexiones de infraestructura (bases de datos, cliente HTTP)
// 2. Parámetros de aplicación
// 3. Repositorios
// 4. Casos de uso
// 5. Servicios
func StartContainer() {
	// 1. Inicializar conexiones de infraestructura
	InitializeConnections()

	// 2. Inicializar parámetros de aplicación
	if err := InitializeApplicationParameters(); err != nil {
		return
	}

	// 3. Inicializar repositorios
	InitializeRepositories()

	// 4. Inicializar casos de uso
	InitializeUseCases()

	// 5. Inicializar servicios
	InitializeServices()
}
