package inject

import (
	"cmp"
	"context"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

var (
	ErrTypeNotProvided     = errors.New("type not provided")
	ErrTypeAlreadyProvided = errors.New("type already provided")
	ErrParentAlreadySet    = errors.New("parent already set")
	ErrTypeNotAllowed      = errors.New("type not allowed")
	ErrCircularDependency  = errors.New("circular dependency detected")
	ErrErrorTypeMustBeLast = errors.New("error type must be the last return value")
	ErrInvalidProvider     = errors.New("provide only accepts a function that returns at least one type except error")
	ErrInvalidInvokeTarget = errors.New("invoke only accepts a function")
	ErrInvalidApplyTarget  = errors.New("apply only accepts a struct")
)

type ctxKeyDependencyPath struct{}

type dependencyPath []reflect.Type

func (dp dependencyPath) String() string {
	if len(dp) == 0 {
		return ""
	}
	var parts []string
	for _, t := range dp {
		parts = append(parts, t.String())
	}
	return strings.Join(parts, " -> ")
}

func getDependencyPath(ctx context.Context) dependencyPath {
	path, _ := ctx.Value(ctxKeyDependencyPath{}).(dependencyPath)
	return path
}

func appendDependencyToPath(ctx context.Context, typ reflect.Type) context.Context {
	path := getDependencyPath(ctx)
	newPath := append(dependencyPath{}, path...)
	newPath = append(newPath, typ)
	return context.WithValue(ctx, ctxKeyDependencyPath{}, newPath)
}

var (
	typeError   = reflect.TypeOf((*error)(nil)).Elem()
	typeContext = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func IsTypeAllowed(typ reflect.Type) bool {
	return typ != typeContext && typ != typeError
}

type provider struct {
	seq  uint64
	ctor any // value func
}

type Injector struct {
	mu sync.RWMutex

	values    map[reflect.Type]reflect.Value
	providers map[reflect.Type]*provider
	parent    *Injector

	maxProviderSeq uint64
	sfg            singleflight.Group
}

func New() *Injector {
	inj := &Injector{
		values:         map[reflect.Type]reflect.Value{},
		providers:      map[reflect.Type]*provider{},
		maxProviderSeq: 0,
	}
	_ = inj.Provide(func() *Injector { return inj })
	return inj
}

func (inj *Injector) SetParent(parent *Injector) error {
	inj.mu.Lock()
	defer inj.mu.Unlock()
	if inj.parent != nil {
		return errors.WithStack(ErrParentAlreadySet)
	}
	inj.parent = parent
	return nil
}

func (inj *Injector) unsafeProvide(ctor any) error {
	rv := reflect.ValueOf(ctor)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return errors.Wrap(ErrInvalidProvider, "ctor is not a function")
	}

	// Get valid output types with error position validation
	outputTypes, err := getValidOutputTypes(rt)
	if err != nil {
		return err
	}

	if len(outputTypes) == 0 {
		return errors.Wrap(ErrInvalidProvider, "no valid output types")
	}

	provider := &provider{
		seq:  inj.maxProviderSeq,
		ctor: ctor,
	}
	inj.maxProviderSeq++

	for _, outType := range outputTypes {
		if _, ok := inj.values[outType]; ok {
			return errors.Wrap(ErrTypeAlreadyProvided, outType.String())
		}

		if _, ok := inj.providers[outType]; ok {
			return errors.Wrap(ErrTypeAlreadyProvided, outType.String())
		}

		inj.providers[outType] = provider
	}

	for _, outType := range outputTypes {
		if err := inj.unsafeDFSDetectCycle(outType); err != nil {
			return err
		}
	}

	return nil
}

func (inj *Injector) invoke(ctx context.Context, f any) ([]reflect.Value, error) {
	rt := reflect.TypeOf(f)
	if rt.Kind() != reflect.Func {
		return nil, errors.Wrap(ErrInvalidInvokeTarget, "f is not a function")
	}

	numIn := rt.NumIn()
	in := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		argType := rt.In(i)
		if argType == typeContext {
			in[i] = reflect.ValueOf(ctx)
			continue
		}
		if !IsTypeAllowed(argType) {
			return nil, errors.Wrap(ErrTypeNotAllowed, argType.String())
		}
		argValue, err := inj.resolve(ctx, argType)
		if err != nil {
			return nil, err
		}
		in[i] = argValue
	}

	outs := reflect.ValueOf(f).Call(in)

	// apply if possible
	for _, out := range outs {
		unwrapped := unwrapPtr(out)
		if unwrapped.Kind() == reflect.Struct {
			if err := inj.applyStruct(ctx, unwrapped); err != nil {
				return nil, err
			}
		}
	}

	numOut := len(outs)
	if numOut > 0 && rt.Out(numOut-1) == typeError {
		rvErr := outs[numOut-1]
		outs = outs[:numOut-1]
		if !rvErr.IsNil() {
			return outs, rvErr.Interface().(error)
		}
	}

	return outs, nil
}

func (inj *Injector) resolve(ctx context.Context, rt reflect.Type) (reflect.Value, error) {
	ctx = appendDependencyToPath(ctx, rt)

	inj.mu.RLock()
	rv := inj.values[rt]
	if rv.IsValid() {
		inj.mu.RUnlock()
		return rv, nil
	}
	provider, ok := inj.providers[rt]
	parent := inj.parent
	inj.mu.RUnlock()

	if ok {
		// ensure that the provider is only executed once same time
		_, err, _ := inj.sfg.Do(strconv.FormatUint(provider.seq, 10), func() (any, error) {
			// must recheck the provider, because it may be deleted by prev inj.sfg.Do
			inj.mu.RLock()
			_, ok := inj.providers[rt]
			inj.mu.RUnlock()
			if !ok {
				return nil, nil
			}

			results, err := inj.invoke(ctx, provider.ctor)
			if err != nil {
				return nil, err
			}

			inj.mu.Lock()
			for _, result := range results {
				resultType := result.Type()
				inj.values[resultType] = result
				delete(inj.providers, resultType)
			}
			inj.mu.Unlock()

			return nil, nil
		})
		if err != nil {
			return rv, err
		}
		return inj.resolve(ctx, rt)
	}

	if parent != nil {
		return parent.resolve(ctx, rt)
	}

	return rv, errors.Wrapf(ErrTypeNotProvided, "dependency path: %s", getDependencyPath(ctx).String())
}

func unwrapPtr(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}

func (inj *Injector) Apply(val any) error {
	return inj.ApplyContext(context.Background(), val)
}

func (inj *Injector) ApplyContext(ctx context.Context, val any) error {
	rv := unwrapPtr(reflect.ValueOf(val))
	if rv.Kind() != reflect.Struct {
		return errors.Wrap(ErrInvalidApplyTarget, "val is not a struct")
	}
	return inj.applyStruct(ctx, rv)
}

var (
	TagName          = "inject"
	TagValueOptional = "optional"
)

func (inj *Injector) applyStruct(ctx context.Context, rv reflect.Value) error {
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		structField := rt.Field(i)
		if tag, ok := structField.Tag.Lookup(TagName); ok {
			if !IsTypeAllowed(structField.Type) {
				return errors.Wrap(ErrTypeNotAllowed, structField.Type.String())
			}

			dep, err := inj.resolve(ctx, structField.Type)
			if err != nil {
				if errors.Is(err, ErrTypeNotProvided) {
					if strings.TrimSpace(tag) == TagValueOptional {
						continue
					}
				}

				return err
			}

			field := rv.Field(i)
			if !field.CanSet() {
				// If the field is unexported, we need to create a new field that is settable.
				field = reflect.NewAt(structField.Type, unsafe.Pointer(field.UnsafeAddr())).Elem()
			}
			field.Set(dep)
		}
	}

	return nil
}

// flatten recursively flattens any arrays/slices found in the arguments.
// This allows grouping related constructors together for better organization.
func flatten(vs ...any) []any {
	var result []any

	for _, f := range vs {
		rv := reflect.ValueOf(f)

		// If it's a slice or array, recursively flatten it
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			// Convert slice/array elements to []any
			subItems := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				subItems[i] = rv.Index(i).Interface()
			}
			// Recursively flatten the sub-items
			result = append(result, flatten(subItems...)...)
		} else {
			// Regular constructor function, add directly
			result = append(result, f)
		}
	}

	return result
}

func (inj *Injector) Provide(ctors ...any) (xerr error) {
	ctors = flatten(ctors...)

	inj.mu.Lock()
	defer inj.mu.Unlock()

	originalProviders := maps.Clone(inj.providers)
	originalMaxSeq := inj.maxProviderSeq
	defer func() {
		if xerr != nil {
			inj.providers = originalProviders
			inj.maxProviderSeq = originalMaxSeq
		}
	}()

	for _, ctor := range ctors {
		if err := inj.unsafeProvide(ctor); err != nil {
			return err
		}
	}

	return nil
}

func (inj *Injector) Invoke(f any) ([]any, error) {
	return inj.InvokeContext(context.Background(), f)
}

func (inj *Injector) InvokeContext(ctx context.Context, f any) ([]any, error) {
	results, err := inj.invoke(ctx, f)
	if err != nil {
		return nil, err
	}
	out := make([]any, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}
	return out, nil
}

func (inj *Injector) Resolve(refs ...any) error {
	return inj.ResolveContext(context.Background(), refs...)
}

func (inj *Injector) ResolveContext(ctx context.Context, refs ...any) error {
	for _, ref := range refs {
		refType := reflect.TypeOf(ref)
		if refType == nil || refType.Kind() != reflect.Ptr {
			return errors.Errorf("resolve requires pointer arguments, got %T", ref)
		}

		rv, err := inj.resolve(ctx, refType.Elem())
		if err != nil {
			return err
		}
		reflect.ValueOf(ref).Elem().Set(rv)
	}
	return nil
}

// Build automatically resolves all provided types using background context.
// This will trigger the creation of all registered constructors,
// ensuring that all dependencies are properly instantiated.
func (inj *Injector) Build(ctors ...any) error {
	return inj.BuildContext(context.Background(), ctors...)
}

// BuildContext automatically resolves all provided types.
// This will trigger the creation of all registered constructors,
// ensuring that all dependencies are properly instantiated.
func (inj *Injector) BuildContext(ctx context.Context, ctors ...any) error {
	if len(ctors) > 0 {
		if err := inj.Provide(ctors...); err != nil {
			return err
		}
	}

	inj.mu.RLock()
	var typesToResolve []reflect.Type
	typeToSeq := make(map[reflect.Type]uint64)
	for typ, provider := range inj.providers {
		typesToResolve = append(typesToResolve, typ)
		typeToSeq[typ] = provider.seq
	}
	inj.mu.RUnlock()

	// Sort types by their provider sequence for stable order based on registration sequence
	slices.SortFunc(typesToResolve, func(a, b reflect.Type) int {
		return cmp.Compare(typeToSeq[a], typeToSeq[b])
	})

	for _, ty := range typesToResolve {
		// Create a pointer to the type for resolution
		ptr := reflect.New(ty)

		// Resolve the type using the injector
		if err := inj.ResolveContext(ctx, ptr.Interface()); err != nil {
			return err
		}
	}

	return nil
}

// getValidOutputTypes extracts all valid output types from a constructor function type,
// performing error type position validation and filtering out error types
func getValidOutputTypes(rt reflect.Type) ([]reflect.Type, error) {
	if rt.Kind() != reflect.Func {
		return nil, nil
	}

	var outputTypes []reflect.Type
	seen := make(map[reflect.Type]bool) // Use map for deduplication
	numOut := rt.NumOut()

	for i := 0; i < numOut; i++ {
		outType := rt.Out(i)

		// Validate error type position: error can only be the last return value
		if outType == typeError {
			if i != numOut-1 {
				return nil, errors.Wrapf(ErrErrorTypeMustBeLast, "error type found at position %d, but must be at position %d", i, numOut-1)
			}
			// Skip error type if it is the last return value
			continue
		}

		if !IsTypeAllowed(outType) {
			return nil, errors.Wrap(ErrTypeNotAllowed, outType.String())
		}

		if !seen[outType] {
			outputTypes = append(outputTypes, outType)
			seen[outType] = true
		}
	}

	return outputTypes, nil
}
