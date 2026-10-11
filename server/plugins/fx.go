package plugins

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/pipeline"
)

var Module = fx.Module(
	"plugins",

	// pipeline собирает реестры inputs/outputs/stages/targetings/limits
	// и менеджер пайплайнов
	pipeline.Module,
)
