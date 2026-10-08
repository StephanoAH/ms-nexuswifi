package infrastructure

type GetSalesRepository interface {
	ExecuteImpl() (string, error)
}
