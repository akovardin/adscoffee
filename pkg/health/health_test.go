package health

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

// mockComponentProvider is a mock implementation of ComponentProvider for testing
type mockComponentProvider struct {
	components []*Component
}

func (m *mockComponentProvider) HealthComponents() []*Component {
	return m.components
}

func TestHealth_Iter(t *testing.T) {
	logger := zaptest.NewLogger(t)

	t.Run("nil callback should not panic", func(t *testing.T) {
		h := &Health{
			logger: logger,
			m:      sync.RWMutex{},
		}

		// This should not panic
		h.Iter(ComponentKindAll, nil)
	})

	t.Run("empty components should call callback zero times", func(t *testing.T) {
		h := &Health{
			logger:     logger,
			components: []*Component{},
			m:          sync.RWMutex{},
		}

		callCount := 0
		h.Iter(ComponentKindAll, func(c *Component) {
			callCount++
		})

		assert.Equal(t, 0, callCount)
	})

	t.Run("components with matching kind should be iterated", func(t *testing.T) {
		component1 := &Component{
			Kind: ComponentKindApp,
			Name: "app-component",
		}
		component2 := &Component{
			Kind: ComponentKindLocal,
			Name: "local-component",
		}
		component3 := &Component{
			Kind: ComponentKindExternal,
			Name: "external-component",
		}

		h := &Health{
			logger:     logger,
			components: []*Component{component1, component2, component3},
			m:          sync.RWMutex{},
		}

		// Test iterating all components
		var iteratedComponents []*Component
		h.Iter(ComponentKindAll, func(c *Component) {
			iteratedComponents = append(iteratedComponents, c)
		})

		assert.Equal(t, 3, len(iteratedComponents))
		assert.Contains(t, iteratedComponents, component1)
		assert.Contains(t, iteratedComponents, component2)
		assert.Contains(t, iteratedComponents, component3)

		// Test iterating only app components
		iteratedComponents = nil
		h.Iter(ComponentKindApp, func(c *Component) {
			iteratedComponents = append(iteratedComponents, c)
		})

		assert.Equal(t, 1, len(iteratedComponents))
		assert.Contains(t, iteratedComponents, component1)

		// Test iterating only local components
		iteratedComponents = nil
		h.Iter(ComponentKindLocal, func(c *Component) {
			iteratedComponents = append(iteratedComponents, c)
		})

		assert.Equal(t, 1, len(iteratedComponents))
		assert.Contains(t, iteratedComponents, component2)

		// Test iterating only external components
		iteratedComponents = nil
		h.Iter(ComponentKindExternal, func(c *Component) {
			iteratedComponents = append(iteratedComponents, c)
		})

		assert.Equal(t, 1, len(iteratedComponents))
		assert.Contains(t, iteratedComponents, component3)
	})

	t.Run("components with non-matching kind should be skipped", func(t *testing.T) {
		component1 := &Component{
			Kind: ComponentKindApp,
			Name: "app-component",
		}
		component2 := &Component{
			Kind: ComponentKindLocal,
			Name: "local-component",
		}

		h := &Health{
			logger:     logger,
			components: []*Component{component1, component2},
			m:          sync.RWMutex{},
		}

		// Test iterating only external components (none should match)
		callCount := 0
		h.Iter(ComponentKindExternal, func(c *Component) {
			callCount++
		})

		assert.Equal(t, 0, callCount)
	})

	t.Run("components from providers should be iterated", func(t *testing.T) {
		component1 := &Component{
			Kind: ComponentKindApp,
			Name: "app-component",
		}
		component2 := &Component{
			Kind: ComponentKindLocal,
			Name: "local-component",
		}

		provider1 := &mockComponentProvider{
			components: []*Component{component1},
		}
		provider2 := &mockComponentProvider{
			components: []*Component{component2},
		}

		h := &Health{
			logger:             logger,
			componentProviders: []ComponentProvider{provider1, provider2},
			m:                  sync.RWMutex{},
		}

		// Test iterating all components from providers
		var iteratedComponents []*Component
		h.Iter(ComponentKindAll, func(c *Component) {
			iteratedComponents = append(iteratedComponents, c)
		})

		assert.Equal(t, 2, len(iteratedComponents))
		assert.Contains(t, iteratedComponents, component1)
		assert.Contains(t, iteratedComponents, component2)
	})

	t.Run("components from both health and providers should be iterated", func(t *testing.T) {
		healthComponent := &Component{
			Kind: ComponentKindApp,
			Name: "health-component",
		}
		providerComponent := &Component{
			Kind: ComponentKindLocal,
			Name: "provider-component",
		}

		provider := &mockComponentProvider{
			components: []*Component{providerComponent},
		}

		h := &Health{
			logger:             logger,
			components:         []*Component{healthComponent},
			componentProviders: []ComponentProvider{provider},
			m:                  sync.RWMutex{},
		}

		// Test iterating all components from both health and providers
		var iteratedComponents []*Component
		h.Iter(ComponentKindAll, func(c *Component) {
			iteratedComponents = append(iteratedComponents, c)
		})

		assert.Equal(t, 2, len(iteratedComponents))
		assert.Contains(t, iteratedComponents, healthComponent)
		assert.Contains(t, iteratedComponents, providerComponent)
	})
}
