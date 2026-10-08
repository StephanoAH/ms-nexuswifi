package container

import (
	"genesis/pos/application/service_application"
	service_application2 "genesis/pos/domain/ports/application/service_application"
)

// GetCurrentTimeService representa el servicio para obtener la hora actual
var GetCurrentTimeService service_application2.GetCurrentTimeService

// ResolveGetCurrentTimeService resuelve el servicio de hora actual (patrón Singleton)
func ResolveGetCurrentTimeService() service_application2.GetCurrentTimeService {
	if GetCurrentTimeService != nil {
		return GetCurrentTimeService
	}
	return &service_application.GetCurrentTimeServiceImpl{
		GetCurrentTimeServiceUseCase: GetCurrentTimeUseCase,
		GetSalesUseCase:              GetSaleUseCase,
	}
}

// InitializeServices inicializa todos los servicios
func InitializeServices() {
	GetCurrentTimeService = ResolveGetCurrentTimeService()
}

