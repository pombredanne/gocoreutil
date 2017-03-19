package main

import (
	"fmt"
	"os"
	"flag"
)

type Pwd struct {
	followSymlink bool
}

func (command *Pwd) addFlag() {
	flag.BoolVar(&command.followSymlink, "L", false, "")
	flag.BoolVar(&command.followSymlink, "P", true, "")
}

func (command *Pwd) Execute() int {
	var pwd string
	var path_contain_symlink bool

	command.addFlag()
	flag.Parse()

	pwd = os.Getenv("PWD")

	path_contain_symlink = false

	if command.followSymlink || path_contain_symlink  {
		working_directory, err := os.Getwd()

		if err != nil {
			return 1
		}

		fmt.Printf("%s\n", working_directory)

	} else {
		fmt.Printf("%s\n", pwd)
	}
	return 0
}
