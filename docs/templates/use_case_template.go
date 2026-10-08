package templates

import (
    "genesis/pos/domain/ports/infrastructure"
)

// [Nombre]UseCaseImpl implementa el caso de uso [Descripción]
type [Nombre]UseCaseImpl struct {
    [Nombre]Repo infrastructure.[Nombre]Repository
}

// ExecuteUseCase ejecuta el caso de uso [Descripción]
func ([variable] *[Nombre]UseCaseImpl) ExecuteUseCase() ([tipo_retorno], error) {
    [resultado], err := [variable].[Nombre]Repo.ExecuteImpl()
    if err != nil {
        return [valor_cero], err
    }
    return [resultado], nil
}

