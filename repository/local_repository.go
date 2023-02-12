package repository

import (
	"os"
	"path/filepath"
	"runtime"
)

const libraryFileName = "lib.so"

type localRepository struct {
	workingDirectory string
	arch             string
	flags            string
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewLocalRepository(workingDirectory string) *localRepository {
	return &localRepository{
		workingDirectory: workingDirectory,
		arch:             runtime.GOARCH,
		flags:            os.Getenv("COMPILE_FLAGS"),
	}
}

func (l *localRepository) Resolve(
	packageName string,
	version string) (found bool, path string, err error) {

	flags := l.flags
	if flags == "" {
		flags = "default"
	}

	path = filepath.Join(l.workingDirectory, packageName, version, l.arch, flags, libraryFileName)

	stat, errTmp := os.Stat(path)
	if errTmp != nil || stat.IsDir() {
		return
	}

	found = true
	return
}
