package repositories_impl

import (
	"context"
	"fmt"
	"genesis/pos/domain/ports/infrastructure"
	"log"
)

type ChargeInitialParamsWithPostgres struct {
	DB infrastructure.DatabaseConnectionPgxInterface
}

func (chargeParameters *ChargeInitialParamsWithPostgres) ChargeInitialParameters() (string, error) {
	type initialParameters struct {
		HostServer string
	}

	parameters := &initialParameters{}

	query := `SELECT VALOR FROM WACHER_PARAMETROS WHERE CODIGO='HOST_SERVER' `
	dbConn, err := chargeParameters.DB.GetDatabaseConnectionPgx()
	if err != nil {
		log.Println("Error connecting to database ::: ", err)
		return "", err
	}
	conn := dbConn.UsePgxDatabaseConnectionDriver()
	defer conn.Close(context.Background())

	resultGet, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Print("ERROR", err)
	}
	fmt.Print("RESULTADO", resultGet)

	log.Println("SendInvoiceElectronic : ", parameters)
	return parameters.HostServer, nil
}
