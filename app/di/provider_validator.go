package di

import (
	"time_management_slackapp/app/validator"

	"github.com/google/wire"
)

func provideContentValidator() validator.ContentValidator{
	return validator.NewContentValidator()
}

var ValidatorSet = wire.NewSet(
	provideContentValidator,
)
