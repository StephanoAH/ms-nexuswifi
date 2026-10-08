package container

import (
	"genesis/pos/domain/ports/infrastructure"
	"genesis/pos/infrastructure/repositories_impl"
)

// GetEnvironmentRepository representa el repositorio de entorno
var GetEnvironmentRepository infrastructure.GetEnvironmentRepository

// ResolveGetEnvironmentRepository resuelve el repositorio de entorno (patrón Singleton)
func ResolveGetEnvironmentRepository() infrastructure.GetEnvironmentRepository {
	if GetEnvironmentRepository != nil {
		return GetEnvironmentRepository
	}
	return &repositories_impl.GetEnvironmentImpl{}
}

// GetCurrentTimeRepository representa el repositorio para obtener la hora actual
var GetCurrentTimeRepository infrastructure.GetCurrentTimeRepository

// ResolveGetCurrentTimeRepository resuelve el repositorio de hora actual (patrón Singleton)
func ResolveGetCurrentTimeRepository() infrastructure.GetCurrentTimeRepository {
	if GetCurrentTimeRepository != nil {
		return GetCurrentTimeRepository
	}
	return &repositories_impl.GetCurrentTimeWithPostgresImpl{
		DB: DatabaseConnectionToLecWithPgx,
	}
}

// GetSaleRepository representa el repositorio de ventas
var GetSaleRepository infrastructure.GetSalesRepository

// ResolveGetSaleRepository resuelve el repositorio de ventas (patrón Singleton)
func ResolveGetSaleRepository() infrastructure.GetSalesRepository {
	if GetSaleRepository != nil {
		return GetSaleRepository
	}
	return &repositories_impl.GetSalesWithMongoImpl{
		MongoDB: DatabaseConnectionToMongo,
	}
}

// BasicAuthRepository representa el repositorio de autenticación básica
var BasicAuthRepository infrastructure.BasicAuthRepositoryPort

// ResolveBasicAuthRepository resuelve el repositorio de autenticación básica (patrón Singleton)
func ResolveBasicAuthRepository() infrastructure.BasicAuthRepositoryPort {
	if BasicAuthRepository != nil {
		return BasicAuthRepository
	}
	BasicAuthRepository = &repositories_impl.BasicAuthRepositoryImpl{
		DB: DatabaseConnectionToLecWithPgx,
	}
	return BasicAuthRepository
}

// InitializeRepositories inicializa todos los repositorios
func InitializeRepositories() {
	GetEnvironmentRepository = ResolveGetEnvironmentRepository()
	GetCurrentTimeRepository = ResolveGetCurrentTimeRepository()
	GetSaleRepository = ResolveGetSaleRepository()
	BasicAuthRepository = ResolveBasicAuthRepository()
}

