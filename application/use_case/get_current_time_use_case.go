package use_case

import (
	"genesis/pos/domain/ports/infrastructure"
	"time"
)

type GetCurrentTimeUseCaseImpl struct {
	GetCurrentTimeRepo infrastructure.GetCurrentTimeRepository
}

func (g *GetCurrentTimeUseCaseImpl) ExecuteUseCase() (time.Time, error) {
	CurrentTime, err := g.GetCurrentTimeRepo.ExecuteImpl()
	return CurrentTime, err
}
