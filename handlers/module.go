package handlers

import (
	"go.uber.org/fx"

	"verifymy-golang-test/services"
)

var Module = fx.Provide(
	services.NewAuthService,
	services.NewUserService,
)
