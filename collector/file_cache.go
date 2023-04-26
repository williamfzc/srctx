package collector

import (
	"os"
	"strings"
)

type FileContentCache struct {
	data map[string][]string
}

func (c *FileContentCache) Get(fileName string) []string {
	if item, ok := c.data[fileName]; ok {
		return item
	}
	// read the file and cache it
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}
	// convert bytes into lines
	lines := strings.Split(string(data), "\n")
	fileCache.data[fileName] = lines
	return fileCache.data[fileName]
}

func (c *FileContentCache) GetLine(fileName string, lineNumber int) string {
	lines := c.Get(fileName)
	if lines == nil {
		return ""
	}
	lineIndex := lineNumber - 1
	return lines[lineIndex]
}

var fileCache = &FileContentCache{data: make(map[string][]string)}
