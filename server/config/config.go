package config

type Config struct {
	Pipelines []Pipeline `yaml:"pipelines"`
}

// тут используется только конфигурация для упорядочивания
// обхода по плагинам
type Pipeline struct {
	Name       string      `yaml:"name"`
	Input      Input       `yaml:"input"`
	Stages     []Stage     `yaml:"stages"`
	Targetings []Targeting `yaml:"targetings"`
	Output     Output      `yaml:"output"`
}

type Input struct {
	Name   string         `yaml:"name"`
	Config map[string]any `yaml:"config"`
}

type Stage struct {
	Name   string         `yaml:"name"`
	Config map[string]any `yaml:"config"`
}

type Targeting struct {
	Name   string         `yaml:"name"`
	Config map[string]any `yaml:"config"`
}

type Output struct {
	Name   string         `yaml:"name"`
	Config map[string]any `yaml:"config"`
}
