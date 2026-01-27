package plugins

type Targeting interface {
	Name() string
	Copy(cfg map[string]any) Targeting
	Filter()
}
