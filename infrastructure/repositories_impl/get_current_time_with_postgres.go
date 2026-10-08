package repositories_impl

import (
	"context"
	"genesis/pos/domain/ports/infrastructure"
	"log"
	"time"
)

type GetCurrentTimeWithPostgresImpl struct {
	DB infrastructure.DatabaseConnectionPgxInterface
}

func (getTime *GetCurrentTimeWithPostgresImpl) ExecuteImpl() (time.Time, error) {
	type timeStruct struct {
		CurrentTime time.Time `json:"CurrentTime"`
	}

	timeInstance := &timeStruct{}
	query := `SELECT NOW() AS "CurrentTime";`

	conn, err := getTime.DB.GetDatabaseConnectionPgx()
	if err != nil {
		log.Println("Error ::: ", err)
		return time.Time{}, err
	}
	defer conn.UsePgxDatabaseConnectionDriver().Close(context.Background())

	resultCurrentTime := conn.UsePgxDatabaseConnectionDriver().QueryRow(context.Background(), query)

	err = resultCurrentTime.Scan(&timeInstance.CurrentTime)

	if err != nil {
		log.Println("Error ::: ", err)
		return time.Time{}, err
	}

	return timeInstance.CurrentTime, nil
}
