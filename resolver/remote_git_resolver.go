package resolver

import (
	"fmt"
	"github.com/hashicorp/go-getter"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"teut.inc/process-engine/compiler"
	"teut.inc/process-engine/repository"
)

var extractorRegex = regexp.MustCompile("^(.*?/.*?/.*?)(/.*)?$")

type remoteGitResolver struct {
	workingDirectory string
	repository       repository.Repository
	compiler         compiler.Compiler
}

func NewRemoteGitResolver(
	workingDirectory string,
	repo repository.Repository,
	comp compiler.Compiler) Resolver {

	if err := os.MkdirAll(workingDirectory, 0750); err != nil {
		log.Fatalf("unable to build working directory for remote git resolver %+v\n", err)
	}

	return &remoteGitResolver{
		workingDirectory: workingDirectory,
		repository:       repo,
		compiler:         comp,
	}
}

func (g *remoteGitResolver) Resolve(addonName string, version string) (found bool, path string, err error) {
	// download the content in our local
	matches := extractorRegex.FindStringSubmatch(addonName)
	repo := matches[1]
	pathInRepo := ""
	if len(matches) > 2 {
		pathInRepo = matches[2][1:]
	}

	var destination string
	destination, err = g.downloadRepository(repo, version)
	if err != nil {
		return
	}

	// compile the addon
	var compiledLibPath string
	compiledLibPath, err = g.compiler.Compile(destination, pathInRepo)
	if err != nil {
		return
	}

	// store the compiled addon in the repository
	path, err = g.repository.Store(compiledLibPath, addonName, version)
	if err == nil {
		found = true
	}
	// fixme: delete compiled lib

	return
}

func (g *remoteGitResolver) downloadRepository(repo string, version string) (destination string, err error) {
	var versionsToTry []string
	if version == "latest" {
		versionsToTry = []string{"main", "master"}
	} else {
		versionsToTry = []string{version}
	}

	destination = filepath.Join(g.workingDirectory, repo)
	err = os.MkdirAll(filepath.Dir(destination), 0750)
	if err != nil {
		return
	}

	errors := make([]any, 0)
	for _, versionToTry := range versionsToTry {
		if err = getter.Get(destination, repo+"?ref="+versionToTry); err == nil {
			return
		} else {
			errors = append(errors, err)
		}
	}

	err = fmt.Errorf(strings.Repeat("%w", len(errors)), errors...)
	return
}
