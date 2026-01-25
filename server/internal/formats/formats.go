package formats

import "go.ads.coffee/server/domain"

type Formats struct {
	list map[string]domain.Format
}

func New(list []domain.Format) *Formats {
	plugins := map[string]domain.Format{}
	for _, format := range list {
		plugins[format.Name()] = format
	}

	return &Formats{
		list: plugins,
	}
}

func (i *Formats) Get(name string, cfg map[string]any) domain.Format {
	return i.list[name].Copy(cfg)
}
