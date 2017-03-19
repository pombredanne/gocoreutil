package main

import (
	"os"
	"path/filepath"
	"fmt"
)

type Dirname struct {
}

func (command *Dirname) Execute() (ret int) {
	if len(os.Args) > 1 {
		fmt.Println(filepath.Dir(os.Args[1]))
		ret = 0
	} else {
		ret = 1
	}
	return 
}
