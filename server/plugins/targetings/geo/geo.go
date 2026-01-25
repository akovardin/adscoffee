package geo

type Geo struct {
}

func New() *Geo {
	return &Geo{}
}

func (g *Geo) Name() string {
	return "geo"
}

func (g *Geo) Filter() {

}
