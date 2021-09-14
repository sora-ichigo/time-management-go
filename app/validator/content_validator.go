package validator

import (
	"starter-restapi-golang/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ContentValidator interface {
	Set(content models.Content) ContentValidator
	Valid() error
}

type contentValidator models.Content

func NewContentValidator() ContentValidator {
	return contentValidator{}
}

func (c contentValidator) Set(content models.Content) ContentValidator {
	c.ID = content.ID
	c.Name = content.Name
	c.Text = content.Text
	return c
}

func (c contentValidator) Valid() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Text, validation.Required),
	)
}
