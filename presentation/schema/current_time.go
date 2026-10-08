package schema

import validation "github.com/go-ozzo/ozzo-validation"

type CurrentTimeBody struct {
	TimeZone string `json:"time_zone"`
	Date     string `json:"date"`
}

func (c CurrentTimeBody) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.TimeZone, validation.Required),
		validation.Field(&c.Date, validation.Date("2006-01-02"), validation.Required),
	)
}

type CurrentTimeParams struct {
	Id int64 `uri:"id"`
}

func (c CurrentTimeParams) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Id, validation.Required),
	)
}

type CurrentTimeQuery struct {
	Authorization string `form:"Authorization"`
}

func (c CurrentTimeQuery) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Authorization, validation.Length(2, 15), validation.Required),
	)
}
