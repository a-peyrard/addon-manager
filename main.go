package main

import (
	"fmt"
	"github.com/a-peyrard/addon-manager/compiler"
	"github.com/a-peyrard/addon-manager/loader"
	"github.com/a-peyrard/addon-manager/process"
	"github.com/a-peyrard/addon-manager/repository"
	"github.com/a-peyrard/addon-manager/resolver"
	"github.com/go-errors/errors"
	"log"
	"os"
	"path/filepath"
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
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get current directory %+v", err)
	}

	localRepository := repository.NewLocalRepository("repo.private")
	defaultCompiler := compiler.NewDefaultCompiler("/tmp/addon/buildOutput")
	processLoader := loader.NewLoader[process.Process](&loader.Config{
		Resolver: resolver.NewAnyResolver([]resolver.Resolver{
			localRepository,
			resolver.NewWorkspaceResolver(
				filepath.Clean(filepath.Join(currentDir, "..")),
				localRepository,
				defaultCompiler,
			),
			resolver.NewRemoteGitResolver(
				"/tmp/addon/gitResolver",
				localRepository,
				defaultCompiler,
			),
		}),
		FactoryMethod: "NewProcess",
	})

	proc, err := processLoader.Load(addonName, version)
	if err != nil {
		log.Fatalf(err.(*errors.Error).ErrorStack())
	}

	if err := proc.Run(); err != nil {
		log.Fatalf(err.(*errors.Error).ErrorStack())
	}
}
