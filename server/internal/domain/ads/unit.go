package ads

type Unit struct {
	ID uint

	Name   string
	Price  int
	Data   string
	Format string

	NetworkID uint
	Network   Network

	PlacementID uint
	Placement   Placement

	Active bool
}
