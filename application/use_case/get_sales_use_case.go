package use_case

import (
	"genesis/pos/domain/ports/infrastructure"
	"log"
)

type GetSalesUseCaseImpl struct {
	GetSalesUseCaseRepo infrastructure.GetSalesRepository
}

func (g *GetSalesUseCaseImpl) ExecuteUseCase() (string, error) {
	r, err := g.GetSalesUseCaseRepo.ExecuteImpl()
	if err != nil {
		log.Println(err)
	}
	return r, nil
}
