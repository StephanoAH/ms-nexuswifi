package service_application

import (
	"genesis/pos/domain/entities"
	"genesis/pos/domain/ports/application/use_case"
	"log"
)

type GetCurrentTimeServiceImpl struct {
	GetCurrentTimeServiceUseCase use_case.GetCurrentTimeUseCase
	GetSalesUseCase              use_case.GetSalesUseCase
}

func (g *GetCurrentTimeServiceImpl) ExecuteService() entities.ResponseGetCurrentTime {
	getCurrentTime, err := g.GetCurrentTimeServiceUseCase.ExecuteUseCase()

	response := &entities.ResponseGetCurrentTime{}

	if err != nil {
		response = &entities.ResponseGetCurrentTime{
			StatusCode: 500,
			Message:    "Error querying the current date",
			Result:     getCurrentTime,
		}
		log.Println("Error ::: ", err)
		return *response
	}

	response = &entities.ResponseGetCurrentTime{
		StatusCode: 200,
		Message:    "The current date is",
		Result:     getCurrentTime,
	}

	return *response
}
