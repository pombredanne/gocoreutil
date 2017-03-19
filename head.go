package main

import (
	"flag"
	"fmt"
	"os"
	"io"
	"bufio"
)

type Head struct {
	versionFlag bool
	numberOfLine int
	printHeaderFlag bool
}

func (command *Head) addFlags() {
	flag.BoolVar(&command.versionFlag, "v", false, "Print the version number.")
	flag.IntVar(&command.numberOfLine, "n", 10, "number of lines")
}

func (command *Head) getHead(file *os.File) (err error) {
	reader := bufio.NewReader(file)
	
	for n := 0; n < command.numberOfLine; n++ {
		line, read_err := reader.ReadString('\n')
		if read_err == io.EOF {
			break
		}	
		if read_err != nil {
			err = read_err
			break
		}
		fmt.Fprintf(os.Stdout, "%s", line)
	}

	return
}

func (command *Head) getFileHead(filename string) (err error) {

	if command.printHeaderFlag {
		fmt.Printf("\n== %s ==\n", filename)	
	}
	
	file, err := os.Open(filename)
	if err != nil {
		return 
	}

	err = command.getHead(file)

	return
}

func (command *Head) getStdin() (err error) {
	err = command.getHead(os.Stdin)
	return	
}

func (command *Head) Execute() (ret int) {
	command.addFlags()
	flag.Parse()

	ret = 0

	if command.versionFlag {
		return ret
	}

	if len(flag.Args()) > 1 {
		command.printHeaderFlag = true
	} else {
		command.printHeaderFlag = false
	}
	
	if len(flag.Args()) == 0 {
		// Head hear the stdin
		err := command.getStdin()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			ret = 1
		}
	} else {
		for _, filename := range flag.Args() {
			err := command.getFileHead(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				ret = 1
			}
		}
	}

	return ret	
}
