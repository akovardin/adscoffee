package formats

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/formats/banner"
	"go.ads.coffee/platform/server/plugins/formats/native"
	"go.ads.coffee/platform/server/plugins/formats/video"
)

var Module = fx.Module(
	"formats.formats",

	banner.Module,
	native.Module,
	video.Module,
)
