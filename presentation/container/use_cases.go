package container

import (
	use_case2 "genesis/pos/application/use_case"
	"genesis/pos/domain/ports/application/use_case"
)

// GetCurrentTimeUseCase representa el caso de uso para obtener la hora actual
var GetCurrentTimeUseCase use_case.GetCurrentTimeUseCase

// ResolveGetCurrentTimeUseCase resuelve el caso de uso de hora actual (patrón Singleton)
func ResolveGetCurrentTimeUseCase() use_case.GetCurrentTimeUseCase {
	if GetCurrentTimeUseCase != nil {
		return GetCurrentTimeUseCase
	}
	return &use_case2.GetCurrentTimeUseCaseImpl{
		GetCurrentTimeRepo: GetCurrentTimeRepository,
	}
}

// GetSaleUseCase representa el caso de uso para obtener ventas
var GetSaleUseCase use_case.GetSalesUseCase

// ResolveGetSaleUseCase resuelve el caso de uso de ventas (patrón Singleton)
func ResolveGetSaleUseCase() use_case.GetSalesUseCase {
	if GetSaleUseCase != nil {
		return GetSaleUseCase
	}
	return &use_case2.GetSalesUseCaseImpl{
		GetSalesUseCaseRepo: GetSaleRepository,
	}
}

// GetCredentialsUseCase representa el caso de uso para obtener credenciales
var GetCredentialsUseCase use_case.GetCredentialsUseCasePort

// ResolveGetCredentialsUseCase resuelve el caso de uso de credenciales (patrón Singleton)
func ResolveGetCredentialsUseCase() use_case.GetCredentialsUseCasePort {
	if GetCredentialsUseCase != nil {
		return GetCredentialsUseCase
	}
	GetCredentialsUseCase = &use_case2.GetCredentialsUseCaseImpl{
		Get: BasicAuthRepository,
	}
	return GetCredentialsUseCase
}

// InitializeUseCases inicializa todos los casos de uso
func InitializeUseCases() {
	GetCurrentTimeUseCase = ResolveGetCurrentTimeUseCase()
	GetSaleUseCase = ResolveGetSaleUseCase()
	GetCredentialsUseCase = ResolveGetCredentialsUseCase()
}

