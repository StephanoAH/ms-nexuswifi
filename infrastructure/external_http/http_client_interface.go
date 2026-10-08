package external_http

type ClientHTTPInterface interface {
	Get(url string, headers map[string]string) ([]byte, error)
	Post(url string, body []byte, headers map[string]string) ([]byte, error)
	Put(url string, body []byte, headers map[string]string) ([]byte, error)
}
