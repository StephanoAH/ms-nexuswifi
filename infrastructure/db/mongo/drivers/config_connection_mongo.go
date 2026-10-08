package drivers

import (
	mongo2 "genesis/pos/infrastructure/db/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
)

type ConfigConnectionMongo struct {
	UrlConnection string
}

func (m *ConfigConnectionMongo) GetConnectionMongo() *mongo2.ConfigDatabaseConnectionMongoDriver {
	mongoOptions := options.Client().ApplyURI(m.UrlConnection)

	client, err := mongo.Connect(mongoOptions)
	if err != nil {
		log.Println("Error mongo Connection ::: ", err)
	}
	return &mongo2.ConfigDatabaseConnectionMongoDriver{MongoClient: client}
}
