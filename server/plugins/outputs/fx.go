package outputs

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/outputs/empty"
	"go.ads.coffee/platform/server/plugins/outputs/pixel"
	"go.ads.coffee/platform/server/plugins/outputs/rtb"
	"go.ads.coffee/platform/server/plugins/outputs/web"
)

var Module = fx.Module(
	"outputs.outputs",

	web.Module,
	rtb.Module,
	pixel.Module,
	empty.Module,
)
