package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	TailModeReadByte = iota
	TailModeReadLine
)

type Tail struct {
	versionFlag     bool
	nrSkipCount     string
	dontTerminate   bool
	printTailerFlag bool
	mode            int
}

func (command *Tail) addFlags() {
	flag.BoolVar(&command.versionFlag, "v", false, "Print the version count.")
	flag.StringVar(&command.nrSkipCount, "n", "10", "count of lines")
	flag.StringVar(&command.nrSkipCount, "c", "10", "count of lines")
	flag.BoolVar(&command.dontTerminate, "f", false, "count of bytes")
}

func (command *Tail) ReadLineFromHead(file *os.File, nrskip int64) (err error) {
	reader := bufio.NewReader(file)

	for i := int64(1); ; i++ {
		line, read_err := reader.ReadString('\n')
		if read_err == io.EOF {
			break
		}
		if read_err != nil {
			err = read_err
			break
		}

		// Skip first nrsikp lines
		if i < nrskip {
			continue
		}

		fmt.Print(line)
	}

	return
}

func (command *Tail) ReadLineFromTail(file *os.File, nrLines int64) (err error) {
	buffer := make([]string, nrLines)
	nrHold := int64(0)
	reader := bufio.NewReader(file)

	for {
		line, read_err := reader.ReadString('\n')

		if read_err == io.EOF {
			break
		}

		if read_err != nil {
			err = read_err
			break
		}

		if nrHold == nrLines {
			for i := int64(1); i < nrLines; i++ {
				buffer[i-1] = buffer[i]
			}
			nrHold--
		}
		buffer[nrHold] = line
		nrHold++

	}

	if err == nil {
		for _, s := range buffer {
			fmt.Print(s)
		}
	}

	return
}

func (command *Tail) ReadByte(file *os.File) (err error) {
	buffer := make([]byte, 512)
	
	for {
		nr_bytes, read_err := file.Read(buffer)
		if read_err == io.EOF {
			break
		}
		
		if read_err != nil {
			err = read_err
			break	
		}

		for i := 0; i < nr_bytes; i++ {
			fmt.Printf("%c", buffer[i])
		}
	}

	return
}

func (command *Tail) ReadByteFromHead(file *os.File, nrskip int64) (err error) {
	file.Seek(nrskip, os.SEEK_SET)
	return command.ReadByte(file)
}

func (command *Tail) ReadByteFromTail(file *os.File, nrReverse int64) (err error) {
	file.Seek(nrReverse, os.SEEK_END)
	return command.ReadByte(file)
}

func (command *Tail) Execute() (ret int) {
	var file *os.File
	var err error

	command.addFlags()
	flag.Parse()

	if command.versionFlag {
		return 0
	}

	command.mode = TailModeReadLine
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "c" {
			command.mode = TailModeReadByte
		} else if f.Name == "n" {
			command.mode = TailModeReadLine
		}
	})

	nrskip, err := strconv.ParseInt(command.nrSkipCount, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return -1
	}

	var skiphead bool
	if strings.HasPrefix(command.nrSkipCount, "+") {
		skiphead = true
	} else {
		skiphead = false
	}

	if flag.NArg() == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "'%s' %s", flag.Arg(0), err)
			return 1
		}
	}
	defer file.Close()

	switch command.mode {
	case TailModeReadLine:
		if skiphead {
			err = command.ReadLineFromHead(file, nrskip)
		} else {
			err = command.ReadLineFromTail(file, nrskip)
		}
		break
	case TailModeReadByte:
		if skiphead {
			err = command.ReadByteFromHead(file, nrskip)
		} else {
			err = command.ReadByteFromTail(file, nrskip)
		}
		break
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		ret = 1
	}

	return ret
}
