package resolver

import (
	"github.com/a-peyrard/addon-manager/compiler"
	"github.com/a-peyrard/addon-manager/repository"
	"golang.org/x/mod/modfile"
	"io/fs"
	"os"
	"path/filepath"
)

type workspaceResolver struct {
	rootDirectory string
	projects      map[string]string
}

func (w *workspaceResolver) findRepo(repo string, version string) (found bool, path string, err error) {
	if version != "latest" {
		return
	}
	path, found = w.projects[repo]
	return
}

func NewWorkspaceResolver(
	rootDirectory string,
	repo repository.Repository,
	comp compiler.Compiler) Resolver {

	return &dynamicResolver{
		sourceManager: &workspaceResolver{
			rootDirectory: rootDirectory,
			projects:      findProjects(rootDirectory, map[string]string{}),
		},
		repository: repo,
		compiler:   comp,
	}
}

func findProjects(dir string, projects map[string]string) map[string]string {
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		return projects
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return projects
	}

	var (
		modPath    string
		parsedFile *modfile.File
		content    []byte
	)
	for _, e := range entries {
		if e.IsDir() || e.Type() == fs.ModeSymlink {
			modPath = filepath.Join(dir, e.Name(), "go.mod")
			stat, err = os.Stat(modPath)
			if os.IsNotExist(err) {
				findProjects(filepath.Join(dir, e.Name()), projects)
			} else if err == nil && stat.Mode().IsRegular() {
				content, err = os.ReadFile(modPath)
				parsedFile, err = modfile.Parse(modPath, content, nil)
				if err != nil {
					return projects
				}

				projects[parsedFile.Module.Mod.Path] = filepath.Join(dir, e.Name())
			}
		}
	}

	return projects
}
