package inputs

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/inputs/postback"
	"go.ads.coffee/platform/server/plugins/inputs/rtb"
	"go.ads.coffee/platform/server/plugins/inputs/static"
	"go.ads.coffee/platform/server/plugins/inputs/tracker"
	"go.ads.coffee/platform/server/plugins/inputs/web"
)

var Module = fx.Module(
	"inputs.inputs",

	rtb.Module,
	web.Module,
	postback.Module,
	tracker.Module,
	static.Module,
)
