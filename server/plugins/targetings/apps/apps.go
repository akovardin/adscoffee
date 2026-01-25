package apps

type Apps struct{}

func New() *Apps {
	return &Apps{}
}

func (a *Apps) Name() string {
	return "apps"
}

func (a *Apps) Filter() {}
