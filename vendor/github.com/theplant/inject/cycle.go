package inject

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// unsafeDFSDetectCycle performs DFS to detect circular dependencies starting from a target type
func (inj *Injector) unsafeDFSDetectCycle(targetType reflect.Type) error {
	return inj.unsafeDFSDetectCycleRecursive(
		targetType,
		make(map[reflect.Type]bool),
		[]directDep{
			{
				rt:  targetType,
				src: "",
			},
		},
	)
}

func (inj *Injector) unsafeDFSDetectCycleRecursive(targetType reflect.Type, visited map[reflect.Type]bool, path []directDep) error {
	if visited[targetType] {
		// Found a cycle - find where the cycle starts
		cycleStart := -1
		for i, dep := range path {
			if dep.rt == targetType {
				cycleStart = i
				break
			}
		}
		if cycleStart >= 0 {
			// Extract the cycle portion
			cyclePart := path[cycleStart:]
			if len(cyclePart) > 0 {
				cycleStrings := make([]string, len(cyclePart))
				for i, dep := range cyclePart {
					cycleStrings[i] = dep.String()
				}
				return errors.Wrap(ErrCircularDependency, strings.Join(cycleStrings, " -> "))
			}
		}
	}

	visited[targetType] = true
	defer func() {
		visited[targetType] = false // Backtrack
	}()

	// Get all dependencies for this type (including co-constructor dependencies)
	dependencies, err := inj.unsafeGetDirectDependencies(targetType)
	if err != nil {
		return err
	}

	// Recursively check each dependency
	for _, depInfo := range dependencies {
		if err := inj.unsafeDFSDetectCycleRecursive(depInfo.rt, visited, append(path, depInfo)); err != nil {
			return err
		}
	}

	return nil
}

// directDep represents a dependency with its source information
type directDep struct {
	rt  reflect.Type
	src string // describes where this dependency comes from
}

// String returns a formatted representation of the dependency
func (d directDep) String() string {
	if d.src == "" {
		return d.rt.String()
	}
	return d.rt.String() + "@" + d.src
}

// unsafeGetDirectDependencies returns all dependencies for a given type, including co-constructor dependencies.
// A type's dependencies include:
// 1. Its struct field inject tag dependency types (highest priority)
// 2. Its provider.ctor's In types (except Context)
// 3. Co-constructor dependencies: dependencies of sibling types that could create indirect cycles
func (inj *Injector) unsafeGetDirectDependencies(targetType reflect.Type) ([]directDep, error) {
	var deps []directDep
	seen := make(map[reflect.Type]bool) // Track seen types

	addFieldDep := func(typ reflect.Type) error {
		fieldDeps, err := getStructFieldDependencies(typ)
		if err != nil {
			return err
		}
		for _, depType := range fieldDeps {
			if !seen[depType] {
				deps = append(deps, directDep{
					rt:  depType,
					src: typ.String(), // Direct field injection dependency
				})
				seen[depType] = true
			}
		}
		return nil
	}

	// 1. Add struct field dependencies first (inject tag dependencies) - these have highest priority
	if err := addFieldDep(targetType); err != nil {
		return nil, err
	}

	// Check if this type has a provider
	if provider, ok := inj.providers[targetType]; ok {
		providerType := reflect.TypeOf(provider.ctor)

		// 2. Add input dependencies (function parameters, except Context)
		for i := 0; i < providerType.NumIn(); i++ {
			inType := providerType.In(i)
			if inType != typeContext && !seen[inType] {
				if !IsTypeAllowed(inType) {
					return nil, errors.Wrap(ErrTypeNotAllowed, inType.String())
				}
				deps = append(deps, directDep{
					rt:  inType,
					src: "", // Direct constructor parameter dependency
				})
				seen[inType] = true
			}
		}

		// 3. Add co-constructor dependencies: dependencies from sibling types that have inject tags
		sameConstructorOutputs, _ := getValidOutputTypes(providerType)
		for _, siblingType := range sameConstructorOutputs {
			if siblingType != targetType {
				// Only add sibling type dependencies if the sibling has inject tag dependencies
				if err := addFieldDep(siblingType); err != nil {
					return nil, err
				}
			}
		}
	}

	return deps, nil
}

// getStructFieldDependencies extracts injection dependencies from struct fields
func getStructFieldDependencies(t reflect.Type) ([]reflect.Type, error) {
	var deps []reflect.Type
	seen := make(map[reflect.Type]bool) // Use map for deduplication

	// Dereference pointer types
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if _, ok := field.Tag.Lookup(TagName); ok {
			if !IsTypeAllowed(field.Type) {
				return nil, errors.Wrap(ErrTypeNotAllowed, field.Type.String())
			}
			if !seen[field.Type] {
				deps = append(deps, field.Type)
				seen[field.Type] = true
			}
		}
	}

	return deps, nil
}
