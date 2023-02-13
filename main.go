package main

import (
	"fmt"
	"github.com/a-peyrard/addon-manager/compiler"
	"github.com/a-peyrard/addon-manager/loader"
	"github.com/a-peyrard/addon-manager/process"
	"github.com/a-peyrard/addon-manager/repository"
	"github.com/a-peyrard/addon-manager/resolver"
	"os"
)

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
	localRepository := repository.NewLocalRepository("repo.private")
	processLoader := loader.NewLoader[process.Process](&loader.Config{
		Resolver: resolver.NewAnyResolver([]resolver.Resolver{
			localRepository,
			resolver.NewRemoteGitResolver(
				"/tmp/addon/gitResolver",
				localRepository,
				compiler.NewDefaultCompiler("/tmp/addon/buildOutput"),
			),
		}),
		FactoryMethod: "NewProcess",
	})

	proc, err := processLoader.Load(addonName, version)
	if err != nil {
		panic(err)
	}

	if err := proc.Run(); err != nil {
		panic(err)
	}
}
