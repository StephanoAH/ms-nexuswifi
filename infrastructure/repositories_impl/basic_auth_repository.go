package repositories_impl

import (
	"context"
	"genesis/pos/domain/entities"
	"genesis/pos/domain/ports/infrastructure"
)

type BasicAuthRepositoryImpl struct {
	DB infrastructure.DatabaseConnectionPgxInterface
}

func (b *BasicAuthRepositoryImpl) GetCredentials(integrationID int64, username string) (entities.BasicAuth, error) {
	ctx := context.Background()
	query := "SELECT ipc.identificador_conexion, ipc.clave_conexion FROM homologacion_integraciones.integrador_parametros_conexion ipc WHERE ipc.integrador_id = $1 AND ipc.identificador_conexion = $2;"
	dbConn, err := b.DB.GetDatabaseConnectionPgx()
	if err != nil {
		return entities.BasicAuth{}, err
	}
	conn := dbConn.UsePgxDatabaseConnectionDriver()
	defer conn.Close(ctx)

	var basicAuth entities.BasicAuth
	err = conn.QueryRow(ctx, query, integrationID, username).Scan(&basicAuth.Username, &basicAuth.Password)
	if err != nil {
		return entities.BasicAuth{}, err
	}

	return basicAuth, nil
}
