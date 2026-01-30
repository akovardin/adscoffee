package pipeline

// тут используется только конфигурация для упорядочивания
// обхода по плагинам
type Config struct {
	Name       string      `yaml:"name"`
	Route      string      `yaml:"route"`
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
