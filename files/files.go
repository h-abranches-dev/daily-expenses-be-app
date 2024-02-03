package files

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type FileWrapper struct {
	path  string
	Lines *[]string
	file  *os.File
}

func NewFileWrapper(filePath string) *FileWrapper {
	fw := &FileWrapper{
		path:  filePath,
		Lines: new([]string),
	}
	if err := fw.readFile(); err != nil {
		return nil
	}
	return fw
}

func (fw *FileWrapper) closeFile() error {
	err := fw.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (fw *FileWrapper) readFile() error {
	var err error
	fw.file, err = os.Open(fw.path)
	if err != nil {
		return err
	}

	newLineIdx := 0
	for {
		newChIdx := int64(newLineIdx)
		var sb strings.Builder
		for {
			if _, err = fw.file.Seek(newChIdx, 0); err != nil {
				if err = fw.closeFile(); err != nil {
					return err
				}
				return err
			}

			b1 := make([]byte, 1)
			var n1 int
			n1, err = fw.file.Read(b1)
			if err != nil {
				if err == io.EOF {
					newLineIdx = -1
				} else {
					if err = fw.closeFile(); err != nil {
						return err
					}
					return err
				}
			}

			if newLineIdx == -1 {
				*fw.Lines = append(*fw.Lines, sb.String())
				if err = fw.closeFile(); err != nil {
					return err
				}
				return nil
			}

			if string(b1[:n1]) == "\n" {
				*fw.Lines = append(*fw.Lines, sb.String())
				newLineIdx = int(newChIdx) + 1
				break
			}
			sb.WriteString(string(b1[:n1]))
			newChIdx++
		}
	}
}

func (fw *FileWrapper) writeFile() error {
	var err error
	fw.file, err = os.OpenFile(fw.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	for i := 0; i < len(*fw.Lines)-1; i++ {
		_, err = fw.file.WriteString(fmt.Sprintf("%s\n", (*fw.Lines)[i]))
		if err != nil {
			if err = fw.closeFile(); err != nil {
				return err
			}
			return err
		}
	}

	if err = fw.closeFile(); err != nil {
		return err
	}

	return nil
}

func (fw *FileWrapper) appendToFile(line string) error {
	f, err := os.OpenFile(fw.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s\n", line))
	if err != nil {
		if err = fw.closeFile(); err != nil {
			return err
		}
		return err
	}

	*fw.Lines = append((*fw.Lines)[:len(*fw.Lines)-1], []string{line, ""}...)

	return nil
}

func (fw *FileWrapper) deleteFile() error {
	err := os.Remove(fw.path)
	if err != nil {
		return err
	}

	fw.file = nil

	return nil
}

func (fw *FileWrapper) AppendLine(newLine string) error {
	if err := fw.appendToFile(newLine); err != nil {
		return err
	}

	return nil
}

func (fw *FileWrapper) ReplaceLine(idxLine int, newLine string) error {
	*fw.Lines = append(append((*fw.Lines)[0:idxLine], newLine), (*fw.Lines)[idxLine+1:len(*fw.Lines)]...)

	if err := fw.deleteFile(); err != nil {
		return err
	}

	if err := fw.writeFile(); err != nil {
		return err
	}

	return nil
}

func (fw *FileWrapper) RemoveLine(idxLine int) error {
	*fw.Lines = append((*fw.Lines)[0:idxLine], (*fw.Lines)[idxLine+1:len(*fw.Lines)-1]...)

	if err := fw.deleteFile(); err != nil {
		return err
	}

	if err := fw.writeFile(); err != nil {
		return err
	}

	return nil
}
