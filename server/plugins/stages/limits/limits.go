package limits

type Limits struct{}

func New() *Limits {
	return &Limits{}
}

func (l *Limits) Name() string {
	return "limits"
}
