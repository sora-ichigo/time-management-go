package di

import (
	"starter-restapi-golang/app/validator"

	"github.com/google/wire"
)

func provideContentValidator() validator.ContentValidator{
	return validator.NewContentValidator()
}

var ValidatorSet = wire.NewSet(
	provideContentValidator,
)
