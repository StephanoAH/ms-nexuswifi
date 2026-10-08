package entities

type BasicAuth struct {
	Username string `json:"username" db:"identificador_conexion"`
	Password string `json:"password" db:"clave_conexion"`
}
