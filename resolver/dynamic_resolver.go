package resolver

import (
	"github.com/a-peyrard/addon-manager/compiler"
	"github.com/a-peyrard/addon-manager/repository"
	"regexp"
)

var extractorRegex = regexp.MustCompile("^(.*?/.*?/.*?)(/.*)?$")

type sourceManager interface {
	findRepo(repo string, version string) (found bool, path string, err error)
}

type dynamicResolver struct {
	sourceManager sourceManager
	repository    repository.Repository
	compiler      compiler.Compiler
}

func (d *dynamicResolver) Resolve(addonName string, version string) (found bool, path string, err error) {
	matches := extractorRegex.FindStringSubmatch(addonName)
	repo := matches[1]
	pathInRepo := ""
	if len(matches) > 2 && matches[2] != "" {
		pathInRepo = matches[2][1:]
	}

	// find the source for the addon
	var repoPath string
	found, repoPath, err = d.sourceManager.findRepo(repo, version)
	if err != nil || !found {
		return
	}

	// compile the addon
	var compiledLibPath string
	compiledLibPath, err = d.compiler.Compile(repoPath, pathInRepo)
	if err != nil {
		return
	}

	// store the compiled addon in the repository
	path, err = d.repository.Store(compiledLibPath, addonName, version)
	if err == nil {
		found = true
	}
	// fixme: delete compiled lib

	return
}
