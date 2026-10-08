package external_http

type HttpResponse struct {
	StatusCode int
	Status     string
	Body       []byte
}
