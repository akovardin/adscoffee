package pipeline

import "go.ads.coffee/platform/server/internal/domain/plugins"

// Реестры плагинов всех видов собраны в пакете пайплайна: именно здесь
// из конфига строится цепочка обработки запроса.

type (
	Inputs     = plugins.Registry[plugins.Input]
	Outputs    = plugins.Registry[plugins.Output]
	Stages     = plugins.Registry[plugins.Stage]
	Targetings = plugins.Registry[plugins.Targeting]
	Limits     = plugins.Registry[plugins.Limit]
)

func newInputs(list []plugins.Provider[plugins.Input]) *Inputs {
	return plugins.NewRegistry[plugins.Input]("inputs", list)
}

func newOutputs(list []plugins.Provider[plugins.Output]) *Outputs {
	return plugins.NewRegistry[plugins.Output]("outputs", list)
}

func newStages(list []plugins.Provider[plugins.Stage]) *Stages {
	return plugins.NewRegistry[plugins.Stage]("stages", list)
}

func newTargetings(list []plugins.Provider[plugins.Targeting]) *Targetings {
	return plugins.NewRegistry[plugins.Targeting]("targetings", list)
}

func newLimits(list []plugins.Provider[plugins.Limit]) *Limits {
	return plugins.NewRegistry[plugins.Limit]("limits", list)
}
