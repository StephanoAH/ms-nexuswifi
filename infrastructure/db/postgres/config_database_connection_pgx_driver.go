package postgres

import "github.com/jackc/pgx/v5"

type ConfigDatabaseConnectionPgxDriver struct {
	PgxConn *pgx.Conn
}

func (configDatabaseDriver *ConfigDatabaseConnectionPgxDriver) UsePgxDatabaseConnectionDriver() *pgx.Conn {
	return configDatabaseDriver.PgxConn
}
