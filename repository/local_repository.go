package repository

import (
	"io"
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

func (l *localRepository) Store(tmpAddonPath string, addonName string, version string) (err error) {
	destinationPath := l.generateAddonPath(addonName, version)
	err = os.MkdirAll(destinationPath, 0750)
	if err != nil {
		return
	}
	return copyFile(tmpAddonPath, destinationPath)
}

func (l *localRepository) generateAddonPath(addonName string, version string) string {
	flags := l.flags
	if flags == "" {
		flags = "default"
	}

	return filepath.Join(l.workingDirectory, addonName, version, l.arch, flags, libraryFileName)
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		_ = in.Close()
	}()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		closeErr := out.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
