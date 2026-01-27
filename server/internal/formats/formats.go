package formats

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Formats struct {
	list map[string]plugins.Format
}

func New(list []plugins.Format) *Formats {
	plugins := map[string]plugins.Format{}
	for _, format := range list {
		plugins[format.Name()] = format
	}

	return &Formats{
		list: plugins,
	}
}

func (i *Formats) Get(name string, cfg map[string]any) plugins.Format {
	return i.list[name].Copy(cfg)
}
