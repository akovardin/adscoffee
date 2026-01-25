package targetings

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/targetings/apps"
	"go.ads.coffee/server/plugins/targetings/geo"
)

var Module = fx.Module(
	"targetings.targetings",

	apps.Module,
	geo.Module,
)
