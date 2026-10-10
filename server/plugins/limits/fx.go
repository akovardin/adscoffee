package limits

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/limits/budget"
	"go.ads.coffee/platform/server/plugins/limits/capping"
)

var Module = fx.Module(
	"limits.limits",

	budget.Module,
	capping.Module,
)
