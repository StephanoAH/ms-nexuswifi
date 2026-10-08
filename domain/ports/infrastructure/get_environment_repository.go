package infrastructure

type GetEnvironmentRepository interface {
	GetEnvironmentString(key string) string
	GetEnvironmentInt(key string) (int, error)
	GetEnvironmentBool(key string) (bool, error)
}
