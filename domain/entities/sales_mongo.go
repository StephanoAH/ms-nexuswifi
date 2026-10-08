package entities

type VentaDetallada struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	IDMovimiento       int      `json:"id_movimiento"`
	V                  int      `json:"__v"`
	Ano                int      `json:"ano"`
	Cantidad           float64  `json:"cantidad"`
	Cara               int      `json:"cara"`
	CentroLogistico    string   `json:"centro_logistico"`
	Codigo             string   `json:"codigo"`
	CodigoProducto     int      `json:"codigo_producto"`
	CodigoSapCliente   string   `json:"codigo_sap_cliente"`
	CodigoSapProducto  string   `json:"codigo_sap_producto"`
	ConsecutivoFactura string   `json:"consecutivo_factura"`
	Descuento          float64  `json:"descuento"`
	Dia                int      `json:"dia"`
	EmpresaID          int      `json:"empresa_id"`
	Fecha              string   `json:"fecha"`
	Hora               string   `json:"hora"`
	IDJornada          int      `json:"id_jornada"`
	Isla               int      `json:"isla"`
	Kilometraje        string   `json:"kilometraje"`
	LecturaFinal       int      `json:"lectura_final"`
	LecturaInicial     int      `json:"lectura_inicial"`
	Manguera           int      `json:"manguera"`
	MedioAutorizacion  string   `json:"medio_autorizacion"`
	Mes                int      `json:"mes"`
	MetodoPago         string   `json:"metodo_pago"`
	NombreCliente      string   `json:"nombre_cliente"`
	NombreEstacion     string   `json:"nombre_estacion"`
	NumeroDocumento    string   `json:"numero_documento"`
	Placa              string   `json:"placa"`
	Precio             int      `json:"precio"`
	PrecioDiferencial  *float64 `json:"precio_diferencial"`
	Producto           string   `json:"producto"`
	Promotor           string   `json:"promotor"`
	Surtidor           int      `json:"surtidor"`
	Tanque             int      `json:"tanque"`
	TipoDocumento      string   `json:"tipo_documento"`
	TipoFactura        string   `json:"tipo_factura"`
	TipoNegocio        string   `json:"tipo_negocio"`
	Turno              string   `json:"turno"`
	Unidad             string   `json:"unidad"`
	ValorImpuestoTotal float64  `json:"valor_impuesto_total"`
	ValorVentaTotal    float64  `json:"valor_venta_total"`
}
