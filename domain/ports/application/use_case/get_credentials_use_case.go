package use_case

import "genesis/pos/domain/entities"

type GetCredentialsUseCasePort interface {
	GetCredentialsUseCase(integrationID int64, username string) (entities.BasicAuth, error)
}
