package plugins

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/internal/formats"
	"go.ads.coffee/server/internal/inputs"
	"go.ads.coffee/server/internal/outputs"
	"go.ads.coffee/server/internal/pipeline"
	"go.ads.coffee/server/internal/stages"
	"go.ads.coffee/server/internal/targetings"
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
