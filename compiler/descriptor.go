package compiler

import (
	"github.com/BurntSushi/toml"
	errors2 "github.com/go-errors/errors"
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

type container struct {
	Addon *Descriptor
}

func ParseDescriptor(descriptorPath string) (descriptor *Descriptor, err error) {
	var container container
	var stat os.FileInfo
	stat, err = os.Stat(descriptorPath)
	if err != nil {
		err = errors2.Wrap(err, 1)
		return
	}
	if !stat.Mode().IsRegular() {
		err = errors2.Errorf("descriptor %s is not a regular file", descriptorPath)
		return
	}
	_, err = toml.DecodeFile(descriptorPath, &container)
	if err == nil {
		descriptor = container.Addon
	} else {
		err = errors2.Wrap(err, 1)
	}

	return
}
