package entities

import "time"

type ResponseGetCurrentTime struct {
	StatusCode int       `json:"StatusCode"`
	Message    string    `json:"Message"`
	Result     time.Time `json:"CurrentTime"`
}
