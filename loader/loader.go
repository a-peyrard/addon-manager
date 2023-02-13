package loader

import (
	"fmt"
	"plugin"
	"reflect"
	"teut.inc/process-engine/resolver"
)

type Config struct {
	Resolver      resolver.Resolver
	FactoryMethod string
}

type loader[T any] struct {
	resolver      resolver.Resolver
	factoryMethod string
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewLoader[T any](config *Config) *loader[T] {
	return &loader[T]{
		resolver:      config.Resolver,
		factoryMethod: config.FactoryMethod,
	}
}

func (l *loader[T]) Load(addonName string, version string) (T, error) {
	var nilValue T
	found, path, err := l.resolver.Resolve(addonName, version)
	if err != nil {
		return nilValue, err
	}
	if !found {
		return nilValue, fmt.Errorf("unable to find %s on version %s", addonName, version)
	}
	p, err := plugin.Open(path)
	if err != nil {
		return nilValue, err
	}

	symProcess, err := p.Lookup(l.factoryMethod)
	if err != nil {
		return nilValue, err
	}

	processFactory, ok := symProcess.(func() T)
	if !ok {
		expectedType := reflect.TypeOf((*T)(nil)).Name()
		foundType := reflect.TypeOf(symProcess).Name()
		return nilValue, fmt.Errorf("unexpected type, got %s but was expecting %s", foundType, expectedType)
	}

	return processFactory(), nil
}
