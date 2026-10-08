package use_case

import (
	"genesis/pos/domain/entities"
	"genesis/pos/domain/ports/infrastructure"
)

type GetCredentialsUseCaseImpl struct {
	Get infrastructure.BasicAuthRepositoryPort
}

func (g *GetCredentialsUseCaseImpl) GetCredentialsUseCase(integrationID int64, username string) (entities.BasicAuth, error) {
	credentials, err := g.Get.GetCredentials(integrationID, username)
	if err != nil {
		return entities.BasicAuth{}, err
	}
	return credentials, nil
}
