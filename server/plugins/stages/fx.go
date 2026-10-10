package stages

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/stages/banners"
	"go.ads.coffee/platform/server/plugins/stages/limit"
	"go.ads.coffee/platform/server/plugins/stages/mediation"
	"go.ads.coffee/platform/server/plugins/stages/rotation"
	"go.ads.coffee/platform/server/plugins/stages/targeting"
)

var Module = fx.Module(
	"stages.stages",

	limit.Module,
	targeting.Module,
	rotation.Module,
	banners.Module,
	mediation.Module,
)
