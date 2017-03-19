package main 

import (
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

type Basename struct {
}

func (command *Basename) Execute() int {

	if len(os.Args) > 1 {
		basename := filepath.Base(os.Args[1])
	
		if len(os.Args) > 2 {
			basename = strings.TrimSuffix(basename, os.Args[2])
		}

		fmt.Println(basename)
		return 0
	} else {
		return 1
	}
}

