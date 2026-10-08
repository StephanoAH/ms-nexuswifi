package use_case

import "time"

type GetCurrentTimeUseCase interface {
	ExecuteUseCase() (time.Time, error)
}
