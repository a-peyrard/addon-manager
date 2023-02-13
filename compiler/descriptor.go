package compiler

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var DescriptorFileName = "addon.toml"

type Descriptor struct {
	Name    string
	Version string
	Build   *BuildDescriptor
}

type BuildDescriptor struct {
	Path string
}

func ParseDescriptor(descriptorPath string) (descriptor *Descriptor, err error) {
	var stat os.FileInfo
	stat, err = os.Stat(descriptorPath)
	if err != nil {
		return
	}
	if !stat.Mode().IsRegular() {
		err = fmt.Errorf("descriptor %s is not a regular file", descriptorPath)
		return
	}
	_, err = toml.DecodeFile(descriptorPath, descriptor)

	return
}
