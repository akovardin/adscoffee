package pipeline

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/inputs"
	"go.ads.coffee/platform/server/plugins/limits"
	"go.ads.coffee/platform/server/plugins/outputs"
	"go.ads.coffee/platform/server/plugins/stages"
	"go.ads.coffee/platform/server/plugins/targetings"
)

var Module = fx.Module(
	"pipeline",

	// плагины
	inputs.Module,
	outputs.Module,
	stages.Module,
	targetings.Module,
	limits.Module,

	// реестры плагинов
	fx.Provide(
		fx.Annotate(newInputs, fx.ParamTags(`group:"inputs"`)),
		fx.Annotate(newOutputs, fx.ParamTags(`group:"outputs"`)),
		fx.Annotate(newStages, fx.ParamTags(`group:"stages"`)),
		fx.Annotate(newTargetings, fx.ParamTags(`group:"targetings"`)),
		fx.Annotate(newLimits, fx.ParamTags(`group:"limits"`)),
	),

	// менеджер пайплайнов
	fx.Provide(NewManager),
)
