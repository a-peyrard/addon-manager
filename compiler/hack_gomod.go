package compiler

import (
	"bytes"
	"fmt"
	"github.com/a-peyrard/addon-manager/util/file"
	"golang.org/x/mod/modfile"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func hackGoMod(repoPath string) (err error) {
	var moduleToHack, currentProjectPath string
	moduleToHack, err = getCurrentModuleName()
	if err != nil {
		return
	}
	moduleToHack = strings.TrimSuffix(moduleToHack, "\n")

	currentProjectPath, err = os.Getwd()
	if err != nil {
		return
	}

	modPath := filepath.Join(repoPath, "go.mod")
	modBackupPath := filepath.Join(repoPath, "go.mod.backup")

	var stat os.FileInfo
	stat, err = os.Stat(modPath)
	if err != nil || !stat.Mode().IsRegular() {
		fmt.Printf("repository %s is not using go module, so nothing to hack", repoPath)
		return
	}
	err = file.Copy(modPath, modBackupPath)
	if err != nil {
		return
	}

	// Parse the go.mod file
	var (
		parsedFile *modfile.File
		content    []byte
	)
	content, err = os.ReadFile(modPath)
	parsedFile, err = modfile.Parse(modPath, content, nil)
	if err != nil {
		return
	}

	var version string
	for _, req := range parsedFile.Require {
		if req.Mod.Path == moduleToHack {
			version = req.Mod.Version
			break
		}
	}
	if version == "" {
		err = fmt.Errorf("unable to find the version of the current module %s", moduleToHack)
		return
	}

	// Append a replace directive to the go.mod file
	err = parsedFile.AddReplace(
		moduleToHack,
		version,
		currentProjectPath,
		"",
	)
	if err != nil {
		return
	}

	// Write the changes to the go.mod file
	content, err = parsedFile.Format()
	if err == nil {
		err = os.WriteFile(modPath, content, 0644)
	}

	return
}

func unHackGoMod(repoPath string) (err error) {
	modPath := filepath.Join(repoPath, "go.mod")
	modBackupPath := filepath.Join(repoPath, "go.mod.backup")
	err = file.Copy(modBackupPath, modPath)
	if err == nil {
		err = os.Remove(modBackupPath)
	}
	return
}

func getCurrentModuleName() (moduleName string, err error) {
	// Execute the `go list` command
	cmd := exec.Command("go", "list", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}

	// Extract the module name from the output
	moduleName = out.String()

	return
}
