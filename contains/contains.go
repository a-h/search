package contains

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

// TextInFile checks whether a file contains the text.
func TextInFile(name, text string) (ok bool, bytesRead int64, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()
	return TextInReader(f, text)
}

// TextInReader checks whether a reader contains the text.
func TextInReader(f io.Reader, text string) (ok bool, bytesRead int64, err error) {
	reader := bufio.NewReader(f)
	var buffer bytes.Buffer
	for {
		var l []byte
		var isPrefix bool
		for {
			l, isPrefix, err = reader.ReadLine()
			if err != nil && err != io.EOF {
				return
			}
			buffer.Write(l)
			bytesRead += int64(len(l))
			if !isPrefix {
				break
			}
			if err == io.EOF {
				break
			}
		}
		line := buffer.String()
		if strings.Contains(line, text) {
			ok = true
			return
		}
		buffer.Reset()
		if err == io.EOF {
			ok = false
			err = nil
			return
		}
	}
}
