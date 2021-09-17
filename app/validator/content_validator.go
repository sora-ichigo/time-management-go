package validator

import (
	"starter-restapi-golang/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ContentValidator interface {
	Set(content models.Content) ContentValidator
	Valid() error
}

type contentValidatorImpl models.Content

func NewContentValidator() ContentValidator {
	return contentValidatorImpl{}
}

func (c contentValidatorImpl) Set(content models.Content) ContentValidator {
	c.ID = content.ID
	c.Name = content.Name
	c.Text = content.Text
	return c
}

func (c contentValidatorImpl) Valid() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Text, validation.Required),
	)
}
