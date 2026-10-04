package targetings

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/targetings/apps"
	"go.ads.coffee/platform/server/plugins/targetings/geo"
	"go.ads.coffee/platform/server/plugins/targetings/placement"
	"go.ads.coffee/platform/server/plugins/targetings/timetable"
)

var Module = fx.Module(
	"targetings.targetings",

	apps.Module,
	geo.Module,
	placement.Module,
	timetable.Module,
)
