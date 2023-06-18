package services

import (
	"go.uber.org/fx"

	"verifymy-golang-test/repositories"
)

var Module = fx.Provide(
	repositories.NewUserRepository,
)
