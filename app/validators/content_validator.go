package validators

import (
	"starter-restapi-golang/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Content models.Content

func (c Content) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Text, validation.Required),
	)
}
