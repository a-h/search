package contains

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// TextInFile checks whether a file contains the text.
func TextInFile(ctx context.Context, name, text string) (ok bool, bytesRead int64, err error) {
	f, fErr := os.Open(name)
	if fErr != nil {
		err = fmt.Errorf("%v: %v", name, fErr)
		return
	}
	defer f.Close()
	ok, bytesRead, err = TextInReader(ctx, f, text)
	if err != nil {
		err = fmt.Errorf("%v: %v", name, err)
	}
	return
}

// TextInReader checks whether a reader contains the text.
func TextInReader(ctx context.Context, f io.Reader, text string) (ok bool, bytesRead int64, err error) {
	reader := bufio.NewReader(f)
	var buffer bytes.Buffer
	for {
		var l []byte
		var isPrefix bool
		for {
			select {
			case <-ctx.Done():
				// Quit early.
				return
			default:
				// Continue.
			}
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
