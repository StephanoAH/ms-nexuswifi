package infrastructure

import "time"

type GetCurrentTimeRepository interface {
	ExecuteImpl() (time.Time, error)
}
