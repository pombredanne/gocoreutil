package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Wc struct {
	byteCountEnable bool
	lineCountEnable bool
	charCountEnable bool
	wordCountEnable bool
}

func (command *Wc) addFlags() {
	flag.BoolVar(&command.byteCountEnable, "c", false,
		"Write to the standard output the number of bytes in each input file.")
	flag.BoolVar(&command.lineCountEnable, "l", false,
		"Write to the standard output the number of <newline> characters  in  each  input file.")
	flag.BoolVar(&command.charCountEnable, "m", false,
		"Write to the standard output the number of characters in each input file.")
	flag.BoolVar(&command.wordCountEnable, "w", false,
		"Write to the standard output the number of words in each input file.")
}

func (command *Wc) count(file *os.File) (nrNewline, nrByte, nrWord, nrChar int64, err error) {
	file.Seek(0, os.SEEK_SET)

	reader := bufio.NewReader(file)

	nrNewline = 0
	nrByte = 0
	nrWord = 0
	nrChar = 0
	for {

		line, read_err := reader.ReadString('\n')

		if read_err == io.EOF {
			break
		}

		if read_err != nil {
			err = read_err
			break
		}

		nrNewline++

		l := strings.TrimSpace(line)
		for _, w := range strings.Split(l, " ") {
			if len(w) != 0 {
				nrWord++
			}
		}
		
		nrByte += int64(len(line))

		sr := strings.NewReader(line)
		nrChar += int64(sr.Len())
	}

	return
}

func (command *Wc) output(filename string, nrNewline, nrByte, nrWord, nrChar int64) {
	var fields []int64
	
	if command.lineCountEnable {
		fields = append(fields, nrNewline)
	}

	if command.wordCountEnable {
		fields = append(fields, nrWord)
	}

	if command.charCountEnable {
		fields = append(fields, nrChar)
	} else if command.byteCountEnable {
		fields = append(fields, nrByte)
	}

	if len(fields) == 0 {
		fields = append(fields, nrNewline)
		fields = append(fields, nrWord)
		fields = append(fields, nrByte)
	}

	for _, v := range fields {
		fmt.Printf("%d " , v)
	}
	fmt.Println(filename)
}

func (command *Wc) Execute() (ret int) {
	flag.Parse()

	if flag.NArg() == 0 {
		return -1
	}

	filenames := flag.Args()

	totalBytes := int64(0)
	totalChars := int64(0)
	totalLines := int64(0)
	totalWords := int64(0)

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			ret = 1
			continue
		}

		nrNewline, nrByte, nrWord, nrChar, err := command.count(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			ret = 2
			continue
		}

		command.output(filename, nrNewline, nrByte, nrWord, nrChar)

		totalBytes += nrByte
		totalChars += nrChar
		totalLines += nrNewline
		totalWords += nrWord
	}

	if len(filenames) > 1 {
		command.output("total", totalLines, totalBytes, totalWords, totalChars)
	}

	return
}
