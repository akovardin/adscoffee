package health

import (
	"context"
	"fmt"
	"time"
)

type ComponentKind uint8

const (
	ComponentKindApp ComponentKind = 1 << iota
	ComponentKindLocal
	ComponentKindExternal

	ComponentKindAll = ComponentKindApp | ComponentKindLocal | ComponentKindExternal
)

func (k *ComponentKind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kindStr string
	if err := unmarshal(&kindStr); err != nil {
		return err
	}

	switch kindStr {
	case "app":
		*k = ComponentKindApp
	case "local":
		*k = ComponentKindLocal
	case "external":
		*k = ComponentKindExternal
	default:
		return fmt.Errorf("invalid component kind: %s", kindStr)
	}

	return nil
}

type CheckFunc func(ctx context.Context) error

type ComponentProvider interface {
	HealthComponents() []*Component
}

type Component struct {
	Kind ComponentKind
	Name string

	CheckFunc     CheckFunc
	CheckErr      error
	CheckDuration time.Duration

	StaticDetails map[Detail]any
}

func (c *Component) Check(ctx context.Context) {
	timeStamp := time.Now()

	c.CheckErr = c.CheckFunc(ctx)
	c.CheckDuration = time.Since(timeStamp)
}
