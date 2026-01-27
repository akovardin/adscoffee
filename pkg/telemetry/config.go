package telemetry

type JaegerConfig struct {
	Enabled       bool    `yaml:"enabled"`
	SamplingRatio float64 `yaml:"sampling"`
	Endpoint      string  `yaml:"endpoint"`
}

type Config struct {
	ServiceName string       `yaml:"-"` // set manually
	Hostname    string       `yaml:"-"` // set manually
	Version     string       `yaml:"-"` // set manually
	Jaeger      JaegerConfig `yaml:"jaeger"`
}
