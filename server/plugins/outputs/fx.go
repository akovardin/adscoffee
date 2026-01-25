package outputs

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/outputs/rtb"
	"go.ads.coffee/server/plugins/outputs/web"
)

var Module = fx.Module(
	"outputs.outputs",

	web.Module,
	rtb.Module,
)
