package repositories

import (
	"go.uber.org/fx"

	"verifymy-golang-test/providers"
)

var Module = fx.Provide(
	providers.NewDBDialector,
	providers.NewDBConnection,
)
