package container

import (
	"genesis/pos/domain/ports/infrastructure"
	"genesis/pos/infrastructure/db/mongo/drivers"
	drivers2 "genesis/pos/infrastructure/db/postgres/drivers"
	"genesis/pos/infrastructure/external_http"
	"genesis/pos/infrastructure/external_http/net_http"
	"genesis/pos/infrastructure/repositories_impl"
	"os"
)

// DatabaseConnectionToMongo representa la conexión a MongoDB
var DatabaseConnectionToMongo infrastructure.DatabaseConnectionMongoInterface

// ResolveDatabaseConnectionToMongo resuelve la conexión a MongoDB (patrón Singleton)
func ResolveDatabaseConnectionToMongo() infrastructure.DatabaseConnectionMongoInterface {
	if DatabaseConnectionToMongo != nil {
		return DatabaseConnectionToMongo
	}
	return &drivers.ConfigConnectionMongo{
		UrlConnection: os.Getenv("CONNECTION_MONGOOSE"),
	}
}

// DatabaseConnectionToLecWithPgx representa la conexión a PostgreSQL
var DatabaseConnectionToLecWithPgx infrastructure.DatabaseConnectionPgxInterface

// ResolveDatabaseConnectionToLecWithPgx resuelve la conexión a PostgreSQL (patrón Singleton)
func ResolveDatabaseConnectionToLecWithPgx() infrastructure.DatabaseConnectionPgxInterface {
	if DatabaseConnectionToLecWithPgx != nil {
		return DatabaseConnectionToLecWithPgx
	}
	return &drivers2.ConfigConnectionPgx{
		UrlToConnect: os.Getenv("CONNECTION_POSTGRESQL"),
	}
}

// ClientHttp representa el cliente HTTP para llamadas externas
var ClientHttp external_http.ClientHTTPInterface

// ResolveClientHttp resuelve el cliente HTTP (patrón Singleton)
func ResolveClientHttp() external_http.ClientHTTPInterface {
	if ClientHttp != nil {
		return ClientHttp
	}
	return &net_http.NetHTTPClient{}
}

// ApplicationParameters contiene los parámetros de configuración de la aplicación
type ApplicationParameters struct {
	HostServer string
}

// AppParams contiene los parámetros de la aplicación
var AppParams ApplicationParameters

// ChargeInitialParamsWithPostgres carga los parámetros iniciales desde PostgreSQL
var ChargeInitialParamsWithPostgres *repositories_impl.ChargeInitialParamsWithPostgres

// ResolveChargeInitialParamsWithPostgres resuelve el repositorio de parámetros iniciales (patrón Singleton)
func ResolveChargeInitialParamsWithPostgres() *repositories_impl.ChargeInitialParamsWithPostgres {
	if ChargeInitialParamsWithPostgres != nil {
		return ChargeInitialParamsWithPostgres
	}
	return &repositories_impl.ChargeInitialParamsWithPostgres{
		DB: DatabaseConnectionToLecWithPgx,
	}
}

// InitializeConnections inicializa todas las conexiones de infraestructura
func InitializeConnections() {
	DatabaseConnectionToMongo = ResolveDatabaseConnectionToMongo()
	DatabaseConnectionToLecWithPgx = ResolveDatabaseConnectionToLecWithPgx()
	ClientHttp = ResolveClientHttp()
}

// InitializeApplicationParameters inicializa los parámetros de la aplicación
func InitializeApplicationParameters() error {
	ChargeInitialParamsWithPostgres = ResolveChargeInitialParamsWithPostgres()
	parameters, err := ChargeInitialParamsWithPostgres.ChargeInitialParameters()
	if err != nil {
		return err
	}
	AppParams.HostServer = parameters
	return nil
}

