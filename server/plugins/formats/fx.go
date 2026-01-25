package formats

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/formats/banner"
	"go.ads.coffee/server/plugins/formats/native"
	"go.ads.coffee/server/plugins/formats/video"
)

var Module = fx.Module(
	"formats.formats",

	banner.Module,
	native.Module,
	video.Module,
)
