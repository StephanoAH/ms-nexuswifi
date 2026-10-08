package infrastructure

import "genesis/pos/domain/entities"

type BasicAuthRepositoryPort interface {
	GetCredentials(integrationID int64, username string) (entities.BasicAuth, error)
}
