package lsif

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// This cache implementation is using a temp file to provide key-value data storage
// It allows to avoid storing intermediate calculations in RAM
// The stored data must be a fixed-size value or a slice of fixed-size values, or a pointer to such data

// 2023/05/21:
// currently we still move it back to RAM
// avoid too much IO harming my disk in development ...

type Cache interface {
	SetEntry(id Id, data interface{}) error
	Entry(id Id, data interface{}) error
	Close() error
	setOffset(id Id) error
	GetReader() io.Reader
}

const (
	CacheTypeFile = "file"
	CacheTypeMem  = "mem"
)

var CacheType = CacheTypeFile

type cacheMem struct {
	file      *File
	chunkSize int64
}

func newCache(filename string, data interface{}) (Cache, error) {
	switch CacheType {
	case CacheTypeMem:
		return newCacheMem(filename, data)
	case CacheTypeFile:
		return newCacheFile(filename, data)
	default:
		return nil, fmt.Errorf("invalid cache type: %v", CacheType)
	}
}

func newCacheMem(filename string, data interface{}) (Cache, error) {
	f := New([]byte{})
	return &cacheMem{file: f, chunkSize: int64(binary.Size(data))}, nil
}

func (c *cacheMem) GetReader() io.Reader {
	return c.file
}

func (c *cacheMem) SetEntry(id Id, data interface{}) error {
	if err := c.setOffset(id); err != nil {
		return err
	}

	return binary.Write(c.file, binary.LittleEndian, data)
}

func (c *cacheMem) Entry(id Id, data interface{}) error {
	if err := c.setOffset(id); err != nil {
		return err
	}

	return binary.Read(c.file, binary.LittleEndian, data)
}

func (c *cacheMem) Close() error {
	// virtual file needs no `close`
	return nil
}

func (c *cacheMem) setOffset(id Id) error {
	offset := int64(id) * c.chunkSize
	_, err := c.file.Seek(offset, io.SeekStart)

	return err
}

type cacheFile struct {
	file      *os.File
	chunkSize int64
}

func (c *cacheFile) GetReader() io.Reader {
	return c.file
}

func newCacheFile(filename string, data interface{}) (Cache, error) {
	f, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, err
	}

	return &cacheFile{file: f, chunkSize: int64(binary.Size(data))}, nil
}

func (c *cacheFile) SetEntry(id Id, data interface{}) error {
	if err := c.setOffset(id); err != nil {
		return err
	}

	return binary.Write(c.file, binary.LittleEndian, data)
}

func (c *cacheFile) Entry(id Id, data interface{}) error {
	if err := c.setOffset(id); err != nil {
		return err
	}

	return binary.Read(c.file, binary.LittleEndian, data)
}

func (c *cacheFile) Close() error {
	return c.file.Close()
}

func (c *cacheFile) setOffset(id Id) error {
	offset := int64(id) * c.chunkSize
	_, err := c.file.Seek(offset, io.SeekStart)

	return err
}
