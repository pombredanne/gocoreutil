package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type MD5Checksum [md5.Size]byte

type MD5 struct {
	Checksum MD5Checksum
	Filename string
}

func CalculateMD5Checksum(filename string) (checksum *MD5, err error) {
	m := new(MD5)

	binary, err := ioutil.ReadFile(filename)
	if err == nil {
		m.Checksum = md5.Sum(binary)
		return m, nil
	} else {
		return nil, err
	}
}

func NewMD5FromString(line string) (checksum *MD5, err error) {
	m := new(MD5)
	string_buffer := make([]string, 2)
	nr_string := 0
	pos := 0

	buf := make([]byte, 2)
	for _, c := range strings.Split(line, "") {
		string_buffer[nr_string] = c
		nr_string++

		if nr_string%2 == 0 {
			v, e := strconv.ParseUint(strings.Join(string_buffer, ""), 16, 8)
			if e != nil {
				err = e
				break
			}

			binary.PutUvarint(buf, v)
			m.Checksum[pos] = buf[0]
			pos++
			nr_string = 0
		}
	}

	if err == nil {
		if pos != md5.Size {
			err = fmt.Errorf("Read line is not md5 binary, size is different expected=%d, actual=%d",
				md5.Size, pos)
		}
	}
	return checksum, err
}


func (m *MD5) Equals(other *MD5) (is_same bool) {
	is_same = true

	if len(m.Checksum) == len(other.Checksum) {
		for i, _ := range m.Checksum {
			if m.Checksum[i] != other.Checksum[i] {
				is_same = false
				break
			}
		}
	} else {
		is_same = false
	}
	return is_same
}

type Md5sum struct {
	helpFlag       bool
	checkmode      bool
	textReadFlag   bool
	binaryReadFlag bool
	bsdStyle       bool
	quiet          bool
	showWarnning   bool
	versionFlag    bool
}

func (command *Md5sum) checkByChecksumFile(filename string) (nr_invalid int, err error) {
	err = nil
	nr_invalid = 0

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nr_invalid, err
	}

	buffer := bytes.NewBuffer(file)
	for {
		line, buf_err := buffer.ReadString('\n')
		if buf_err == io.EOF {
			break
		}
		if buf_err != nil {
			err = buf_err
			break
		}

		// Skip comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)

		var checksum_field string
		var filename_field string

		if len(fields) == 2 {
			checksum_field = fields[0]
			filename_field = fields[1]
		} else if len(fields) == 4 {
			// BSD style checksum file
			if fields[0] != "MD5" {
				continue
			}

			checksum_field = fields[3]
			fmt.Sscanf(fields[1], "(%s)", &filename_field)
		}

		expected_checksum, chk_err := NewMD5FromString(checksum_field)
		if chk_err != nil {
			err = chk_err
			break
		}

		actual_checksum, chk_err := CalculateMD5Checksum(filename_field)
		if chk_err != nil {
			err = chk_err
			break
		}

		if actual_checksum.Equals(expected_checksum) {
			if !command.quiet {
				fmt.Printf("%s: OK\n", filename_field)
			}
		} else {
			nr_invalid++
			fmt.Printf("%s FAILD\n", filename_field)
		}

	}
	return nr_invalid, err
}

func (command *Md5sum) printChecksum(filename string) (err error) {

	md5sum, err := CalculateMD5Checksum(filename)
	if err != nil {
		return err
	}

	if command.bsdStyle {
		fmt.Printf("MD5 (%s) = ", md5sum.Filename)
	}

	for _, b := range md5sum.Checksum {
		fmt.Printf("%02x", b)
	}

	if command.bsdStyle {
		fmt.Println("")
	} else {
		fmt.Printf(" %s\n", md5sum.Filename)
	}
	return nil
}

func (command *Md5sum) addFlag() {
	flag.BoolVar(&command.textReadFlag, "t", false, "print version")
	flag.BoolVar(&command.binaryReadFlag, "b", true, "print version")
	flag.BoolVar(&command.checkmode, "c", false, "check from file")
	flag.BoolVar(&command.bsdStyle, "T", false, "check from file")
	flag.BoolVar(&command.quiet, "q", false, "Do not print 'OK'")
	flag.BoolVar(&command.showWarnning, "-warn", false, "check from file")
	flag.BoolVar(&command.versionFlag, "V", false, "print version")
}


func (command *Md5sum) Execute() int {
	command.addFlag()
	flag.Parse() // Scan the arguments list

	if command.versionFlag {
		fmt.Println("Version:", APP_VERSION)
		return 0
	}

	if command.checkmode {
		nr_invalid := 0
		for _, filename := range flag.Args() {
			n, err := command.checkByChecksumFile(filename)
			if err != nil {
				fmt.Println(err)
				break
			}
			nr_invalid += n
		}

		if nr_invalid != 0 {
			fmt.Printf("md5sum: WARNING: %d computed checksum did not much\n", nr_invalid)
			return 1
		}
	} else {
		for _, filename := range flag.Args() {
			err := command.printChecksum(filename)
			if err != nil {
				fmt.Println(err)
				return 1
			}
		}
	}
	return 0
}
