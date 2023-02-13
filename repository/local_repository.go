package repository

import (
	"os"
	"path/filepath"
	"runtime"
	"teut.inc/process-engine/util/file"
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
	addonName string,
	version string) (found bool, path string, err error) {

	path = l.generateAddonPath(addonName, version)

	stat, errTmp := os.Stat(path)
	if errTmp != nil || stat.IsDir() {
		return
	}

	found = true
	return
}

func (l *localRepository) Store(tmpAddonPath string, addonName string, version string) (path string, err error) {
	path = l.generateAddonPath(addonName, version)
	err = os.MkdirAll(filepath.Dir(path), 0750)
	if err == nil {
		err = file.Copy(tmpAddonPath, path)
	}

	return
}

func (l *localRepository) generateAddonPath(addonName string, version string) string {
	flags := l.flags
	if flags == "" {
		flags = "default"
	}

	return filepath.Join(l.workingDirectory, addonName, version, l.arch, flags, libraryFileName)
}
