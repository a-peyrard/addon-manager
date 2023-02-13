package main

import (
	"fmt"
	"os"
	"plugin"
	"teut.inc/process-engine/loader"
	"teut.inc/process-engine/process"
	"teut.inc/process-engine/repository"
	"teut.inc/process-engine/resolver"
)

func LoadProcess[T any](filepath string, factoryMethod string) (T, error) {
	var nilValue T
	p, err := plugin.Open(filepath)
	if err != nil {
		return nilValue, err
	}

	symProcess, err := p.Lookup(factoryMethod)
	if err != nil {
		return nilValue, err
	}

	processFactory, ok := symProcess.(func() T)
	if !ok {
		return nilValue, fmt.Errorf("unexpected type from module symbol")
	}

	return processFactory(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <package name> (<package version>)\n",
			os.Args[0])
		os.Exit(1)
	}
	addonName := os.Args[1]
	version := "latest"
	if len(os.Args) > 2 {
		version = os.Args[2]
	}
	processLoader := loader.NewLoader[process.Process](&loader.Config{
		Resolver: resolver.NewAnyResolver([]resolver.Resolver{
			repository.NewLocalRepository("repo.private"),
			resolver.NewGithubResolver(),
		}),
		FactoryMethod: "NewProcess",
	})
	proc, err := processLoader.Load(addonName, version)

	//proc, err := LoadProcess[process.Process](os.Args[1], "NewProcess")
	if err != nil {
		panic(err)
	}

	if err := proc.Run(); err != nil {
		panic(err)
	}
}

//Note: The plugins must be built as shared libraries (.so files on Unix-like systems or .dll files on Windows) for the plugin package to be able to load them.
