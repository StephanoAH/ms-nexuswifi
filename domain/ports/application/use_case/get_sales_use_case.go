package use_case

type GetSalesUseCase interface {
	ExecuteUseCase() (string, error)
}
