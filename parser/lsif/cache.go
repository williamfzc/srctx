package lsif

import (
	"encoding/binary"
	"io"
)

// This cache implementation is using a temp file to provide key-value data storage
// It allows to avoid storing intermediate calculations in RAM
// The stored data must be a fixed-size value or a slice of fixed-size values, or a pointer to such data

// 2023/05/21:
// currently we still move it back to RAM
// avoid too much IO harming my disk in development ...
type cache struct {
	file      *File
	chunkSize int64
}

func newCache(filename string, data interface{}) (*cache, error) {
	f := New([]byte{})
	return &cache{file: f, chunkSize: int64(binary.Size(data))}, nil
}

func (c *cache) SetEntry(id Id, data interface{}) error {
	if err := c.setOffset(id); err != nil {
		return err
	}

	return binary.Write(c.file, binary.LittleEndian, data)
}

func (c *cache) Entry(id Id, data interface{}) error {
	if err := c.setOffset(id); err != nil {
		return err
	}

	return binary.Read(c.file, binary.LittleEndian, data)
}

func (c *cache) Close() error {
	// virtual file needs no `close`
	return nil
}

func (c *cache) setOffset(id Id) error {
	offset := int64(id) * c.chunkSize
	_, err := c.file.Seek(offset, io.SeekStart)

	return err
}
