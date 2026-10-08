package templates

import (
    "genesis/pos/domain/entities"
    "genesis/pos/domain/ports/application/use_case"
    "log"
)

// [Nombre]ServiceImpl implementa el servicio [Descripción]
type [Nombre]ServiceImpl struct {
    [Nombre]UseCase use_case.[Nombre]UseCase
    // Agregar más casos de uso si es necesario
    // [Otro]UseCase use_case.[Otro]UseCase
}

// ExecuteService ejecuta el servicio [Descripción]
func ([variable] *[Nombre]ServiceImpl) ExecuteService() entities.Response[Nombre] {
    [resultado], err := [variable].[Nombre]UseCase.ExecuteUseCase()

    response := &entities.Response[Nombre]{}

    if err != nil {
        response = &entities.Response[Nombre]{
            StatusCode: 500,
            Message:    "Error [descripción del error]",
            Result:     [valor_cero],
        }
        log.Println("Error ::: ", err)
        return *response
    }

    response = &entities.Response[Nombre]{
        StatusCode: 200,
        Message:    "[Mensaje de éxito]",
        Result:     [resultado],
    }

    return *response
}

