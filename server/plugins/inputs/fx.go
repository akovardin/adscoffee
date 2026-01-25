package inputs

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/inputs/rtb"
	"go.ads.coffee/server/plugins/inputs/web"
)

var Module = fx.Module(
	"inputs.inputs",

	rtb.Module,
	web.Module,
)
