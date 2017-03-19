package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const APP_VERSION = "0.1"

type CoreutilCommand interface {
	Execute() int
}

var CoreutilCommandTable = map[string]CoreutilCommand{
	"coreutils": new(Coreutils),
	"pwd":       new(Pwd),
	"basename":  new(Basename),
	"md5sum":    new(Md5sum),
	"head":      new(Head),
	"tail":      new(Tail),
	"dirname":   new(Dirname),
	"wc":   new(Wc),
}

type Coreutils struct {
	versionFlag bool
}

func (c *Coreutils) addFlags() {
	// The flag package provides a default help printer via -h switch
	flag.BoolVar(&c.versionFlag, "v", false, "Print the version number.")
}

func (c *Coreutils) Help() {
}

func (c *Coreutils) Execute() int {
	flag.Parse()

	if c.versionFlag {
		fmt.Println("Version:", APP_VERSION)
	} else {
		c.Help()
	}
	return 0
}

func main() {
	command_path := os.Args[0]
	command_name := filepath.Base(command_path)
	command, ok := CoreutilCommandTable[command_name]

	if ok {
		os.Exit(command.Execute())
	} else {
		fmt.Printf("No such command '%s'\n", command_name)
		os.Exit(-1)
	}
}
