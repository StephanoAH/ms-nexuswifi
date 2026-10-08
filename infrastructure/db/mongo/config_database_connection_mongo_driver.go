package mongo

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ConfigDatabaseConnectionMongoDriver struct {
	MongoClient *mongo.Client
}

func (configDatabaseDriver *ConfigDatabaseConnectionMongoDriver) UseMongoDatabaseConnectionDriver() *mongo.Client {
	return configDatabaseDriver.MongoClient
}
