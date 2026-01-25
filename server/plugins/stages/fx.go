package stages

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/stages/banners"
	"go.ads.coffee/server/plugins/stages/limits"
	"go.ads.coffee/server/plugins/stages/rotation"
	"go.ads.coffee/server/plugins/stages/targeting"
)

var Module = fx.Module(
	"stages.stages",

	limits.Module,
	targeting.Module,
	rotation.Module,
	banners.Module,
)
