package plugins

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/pipeline"
	"go.ads.coffee/server/plugins/inputs"
	"go.ads.coffee/server/plugins/inputs/rtb"
	"go.ads.coffee/server/plugins/outputs"
	"go.ads.coffee/server/plugins/outputs/banner"
	"go.ads.coffee/server/plugins/stages"
	"go.ads.coffee/server/plugins/stages/limits"
	"go.ads.coffee/server/plugins/stages/targeting"
	"go.ads.coffee/server/plugins/targetings"
	"go.ads.coffee/server/plugins/targetings/apps"
	"go.ads.coffee/server/plugins/targetings/geo"
)

var Module = fx.Module(
	"plugins",

	// inputs
	inputs.Module,
	rtb.Module,

	// stages
	stages.Module,
	limits.Module,
	targeting.Module,

	// targetings
	targetings.Module,
	apps.Module,
	geo.Module,

	// output
	outputs.Module,
	banner.Module,

	// pipeline
	pipeline.Module,
)
