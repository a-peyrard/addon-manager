package compiler

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
)

type Compiler interface {
	Compile(addonPath string) (compiledLibPath string, err error)
}

type defaultCompiler struct {
	workingDirectory string
	compilationUnit  uint64
	flags            string
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewDefaultCompiler(workingDirectory string) Compiler {
	if err := os.MkdirAll(workingDirectory, 0750); err != nil {
		log.Fatalf("unable to build working directory for default compiler %+v\n", err)
	}

	return &defaultCompiler{
		workingDirectory: workingDirectory,
		compilationUnit:  1,
		flags:            os.Getenv("COMPILE_FLAGS"),
	}
}

func (d *defaultCompiler) Compile(addonPath string) (compiledLibPath string, err error) {
	// get descriptor
	var descriptor *Descriptor
	descriptor, err = ParseDescriptor(filepath.Join(addonPath, DescriptorFileName))
	if err != nil {
		return
	}

	// compile
	compiledLibPath = filepath.Join(
		d.workingDirectory,
		// in order to be thread safe, always generate a new output name
		fmt.Sprintf("addon%d.so", atomic.AddUint64(&d.compilationUnit, 1)),
	)

	commandArgs := make([]string, 0)
	commandArgs = append(commandArgs, "build")
	if d.flags != "" {
		commandArgs = append(commandArgs, fmt.Sprintf("-gcflags=\"%s\"", d.flags))
	}
	commandArgs = append(commandArgs, "-buildmode=plugin")
	commandArgs = append(commandArgs, "-o")
	commandArgs = append(commandArgs, compiledLibPath)
	commandArgs = append(commandArgs, d.generateBuildPath(addonPath, descriptor))

	cmd := exec.Command("go", commandArgs...)
	err = cmd.Run()

	return
}

func (d *defaultCompiler) generateBuildPath(addonPath string, descriptor *Descriptor) (buildPath string) {
	if descriptor.Build == nil || descriptor.Build.Path == "" {
		buildPath = filepath.Join(addonPath, "*.go")
	} else {
		buildPath = descriptor.Build.Path
	}

	return
}
