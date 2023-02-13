package compiler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type Compiler interface {
	Compile(repoPath, pathInRepo string) (compiledLibPath string, err error)
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
		compilationUnit:  0,
		flags:            os.Getenv("COMPILE_FLAGS"),
	}
}

func (d *defaultCompiler) Compile(repoPath, pathInRepo string) (compiledLibPath string, err error) {
	// get descriptor
	var descriptor *Descriptor
	descriptor, err = ParseDescriptor(filepath.Join(repoPath, pathInRepo, DescriptorFileName))
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
		commandArgs = append(commandArgs, fmt.Sprintf("-gcflags=%s", d.flags))
	}
	commandArgs = append(commandArgs, "-buildmode=plugin")
	commandArgs = append(commandArgs, "-o")
	commandArgs = append(commandArgs, compiledLibPath)
	commandArgs = append(commandArgs, generateBuildPath(repoPath, pathInRepo, descriptor))

	err = execCompilationCommand(repoPath, commandArgs)

	if err != nil {
		err = fmt.Errorf(
			"unable to compile the library, command used was '%s %s'",
			"go", strings.Join(commandArgs, " "),
		)
	}

	return
}

func generateBuildPath(repoPath, pathInRepo string, descriptor *Descriptor) (buildPath string) {
	if descriptor.Build == nil || descriptor.Build.Path == "" {
		matches, _ := filepath.Glob(filepath.Join(repoPath, pathInRepo, "*.go"))
		paths := make([]string, len(matches))
		for idx, match := range matches {
			rel, _ := filepath.Rel(repoPath, match)
			paths[idx] = rel
		}
		buildPath = strings.Join(paths, " ")
	} else {
		// fixme we probably want to revisit this, as the files will not have there correct absolute path
		buildPath = descriptor.Build.Path
	}

	return
}

func execCompilationCommand(repoPath string, commandArgs []string) (err error) {
	var currentPath string
	currentPath, err = os.Getwd()
	if err != nil {
		return
	}
	// let's execute the command at the root directory of the addon repository
	if err = os.Chdir(repoPath); err != nil {
		return
	}
	defer func() {
		_ = os.Chdir(currentPath)
	}()

	cmd := exec.Command("go", commandArgs...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err = cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// fixme: make the output better, maybe we shall capture this in a string, and only display it in case of
		// fixme: error
		fmt.Println(scanner.Text())
	}

	// Read the error from the command
	scanner = bufio.NewScanner(stderr)
	for scanner.Scan() {
		// fixme: make the output better
		fmt.Println(scanner.Text())
	}
	err = cmd.Wait()

	return
}
