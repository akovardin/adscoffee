package plugins

import (
	"encoding/json"
	"fmt"
)

// Provider связывает имя плагина с фабрикой, которая создаёт его по конфигу
// пайплайна. Фабрика возвращает ошибку, поэтому некорректный конфиг ловится
// на старте сервера, а не паникует во время обработки запроса.
type Provider[T any] struct {
	Name    string
	Factory func(cfg map[string]any) (T, error)
}

// Registry — реестр фабрик плагинов одного вида (inputs, outputs, stages, …).
type Registry[T any] struct {
	kind  string
	items map[string]func(cfg map[string]any) (T, error)
}

func NewRegistry[T any](kind string, list []Provider[T]) *Registry[T] {
	items := make(map[string]func(cfg map[string]any) (T, error), len(list))
	for _, p := range list {
		items[p.Name] = p.Factory
	}

	return &Registry[T]{
		kind:  kind,
		items: items,
	}
}

func (r *Registry[T]) Get(name string, cfg map[string]any) (T, error) {
	factory, ok := r.items[name]
	if !ok {
		var zero T

		return zero, fmt.Errorf("%s %q is not registered", r.kind, name)
	}

	return factory(cfg)
}

// Decode переносит конфиг вида map[string]any в типизированную структуру
// плагина. Пустой конфиг даёт значения по умолчанию.
func Decode[T any](cfg map[string]any) (T, error) {
	var out T
	if len(cfg) == 0 {
		return out, nil
	}

	raw, err := json.Marshal(cfg)
	if err != nil {
		return out, fmt.Errorf("failed to encode config: %w", err)
	}

	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("failed to decode config: %w", err)
	}

	return out, nil
}
