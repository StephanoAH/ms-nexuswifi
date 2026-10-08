package repositories_impl

import (
	"context"
	"fmt"
	"genesis/pos/domain/entities"
	"genesis/pos/domain/ports/infrastructure"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
)

type GetSalesWithMongoImpl struct {
	MongoDB infrastructure.DatabaseConnectionMongoInterface
}

func (m *GetSalesWithMongoImpl) ExecuteImpl() (string, error) {
	conn := m.MongoDB.GetConnectionMongo().UseMongoDatabaseConnectionDriver()

	collection := conn.Database("afrodita").Collection("venta_detallada")

	result, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error ::: ", err)
	}

	defer func(result *mongo.Cursor, ctx context.Context) {
		err := result.Close(ctx)
		if err != nil {

		}
	}(result, context.Background())

	var values []entities.VentaDetallada
	for result.Next(context.Background()) {
		var value entities.VentaDetallada
		if err := result.Decode(&value); err != nil {
			log.Println("Error decoding document :", err)
			continue
		}
		values = append(values, value)
	}
	fmt.Print("Response ::", values[0])
	return "", err
}
