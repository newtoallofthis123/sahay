package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type File struct {
	// Indexes the file of upto  2^16 lines
	Seek map[uint16]string

	file *os.File
}

func NewFile(file *os.File) (File, error) {
	seek := make(map[uint16]string, 0)
	seek[0] = ""

	reader := bufio.NewReader(file)

	i := 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return File{}, fmt.Errorf("error reading file: %v", err)
		}

		if err == io.EOF {
			break
		}
		seek[uint16(i)] = string(line)
		i++
	}

	return File{
		Seek: seek,
		file: file,
	}, nil
}

func (f *File) GetContents(lineNumber uint16) (string, bool) {
	if _, ok := f.Seek[lineNumber]; !ok {
		return "", false
	}

	return f.Seek[lineNumber], true
}
