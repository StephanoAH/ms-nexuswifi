package service_application

import "genesis/pos/domain/entities"

type GetCurrentTimeService interface {
	ExecuteService() entities.ResponseGetCurrentTime
}
