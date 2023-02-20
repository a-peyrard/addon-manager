package resolver

import (
	"github.com/a-peyrard/addon-manager/compiler"
	"github.com/a-peyrard/addon-manager/repository"
	errors2 "github.com/go-errors/errors"
	"github.com/hashicorp/go-getter"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type remoteGitSourceManager struct {
	workingDirectory string
}

func (r *remoteGitSourceManager) findRepo(repo string, version string) (found bool, path string, err error) {
	var versionsToTry []string
	if version == "latest" {
		versionsToTry = []string{"main", "master"}
	} else {
		versionsToTry = []string{version}
	}

	path = filepath.Join(r.workingDirectory, repo)
	err = os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		err = errors2.Wrap(err, 1)
		return
	}

	// fixme: use go 1.20 errors.Join
	errors := make([]any, 0)
	for _, versionToTry := range versionsToTry {
		if err = getter.Get(path, repo+"?ref="+versionToTry); err == nil {
			found = true
			return
		} else {
			errors = append(errors, errors2.Wrap(err, 1))
		}
	}

	err = errors2.Errorf(strings.Repeat("%w", len(errors)), errors...)
	return
}

func NewRemoteGitResolver(
	workingDirectory string,
	repo repository.Repository,
	comp compiler.Compiler) Resolver {

	if err := os.MkdirAll(workingDirectory, 0750); err != nil {
		log.Fatalf("unable to build working directory for remote git resolver %+v\n", err)
	}

	return &dynamicResolver{
		sourceManager: &remoteGitSourceManager{
			workingDirectory: workingDirectory,
		},
		repository: repo,
		compiler:   comp,
	}
}
