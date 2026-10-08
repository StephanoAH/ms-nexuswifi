package infrastructure

import (
	"genesis/pos/infrastructure/db/mongo"
	"genesis/pos/infrastructure/db/postgres"
)

type DatabaseConnectionMongoInterface interface {
	GetConnectionMongo() *mongo.ConfigDatabaseConnectionMongoDriver
}
type DatabaseConnectionPgxInterface interface {
	GetDatabaseConnectionPgx() (*postgres.ConfigDatabaseConnectionPgxDriver, error)
}
