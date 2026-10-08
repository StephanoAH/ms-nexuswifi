package drivers

import (
	"context"
	"genesis/pos/infrastructure/db/postgres"
	"github.com/jackc/pgx/v5"
	"log"
)

type ConfigConnectionPgx struct {
	UrlToConnect string
}

func (configConnection *ConfigConnectionPgx) GetDatabaseConnectionPgx() (*postgres.ConfigDatabaseConnectionPgxDriver, error) {
	connection, err := pgx.Connect(context.Background(), configConnection.UrlToConnect)
	if err != nil {
		log.Println("Error connecting to database ::: ", err)
		return nil, err
	}
	return &postgres.ConfigDatabaseConnectionPgxDriver{PgxConn: connection}, nil
}
