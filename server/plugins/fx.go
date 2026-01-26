package plugins

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/formats"
	"go.ads.coffee/platform/server/internal/inputs"
	"go.ads.coffee/platform/server/internal/outputs"
	"go.ads.coffee/platform/server/internal/pipeline"
	"go.ads.coffee/platform/server/internal/stages"
	"go.ads.coffee/platform/server/internal/targetings"
)

var Module = fx.Module(
	"plugins",

	// formats
	formats.Module,

	// inputs
	inputs.Module,

	// stages
	stages.Module,

	// targetings
	targetings.Module,

	// output
	outputs.Module,

	// pipeline
	pipeline.Module,
)
