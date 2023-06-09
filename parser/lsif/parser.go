package lsif

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

var Lsif = "lsif"

type Parser struct {
	Docs *Docs

	pr *io.PipeReader
}

func NewParserRaw(ctx context.Context, r io.Reader) (*Parser, error) {
	docs, err := NewDocs()
	if err != nil {
		return nil, err
	}

	if err := docs.Parse(r); err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	parser := &Parser{
		Docs: docs,
		pr:   pr,
	}

	err = pw.Close()
	if err != nil {
		return nil, err
	}

	return parser, nil
}

func NewParser(ctx context.Context, r io.Reader) (*Parser, error) {
	docs, err := NewDocs()
	if err != nil {
		return nil, err
	}

	// ZIP files need to be seekable. Don't hold it all in RAM, use a tempfile
	tempFile, err := os.CreateTemp("", Lsif)
	if err != nil {
		return nil, err
	}

	size, err := io.Copy(tempFile, r)
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(tempFile, size)
	if err != nil {
		return nil, err
	}

	if len(zr.File) == 0 {
		return nil, errors.New("empty zip file")
	}

	file, err := zr.File[0].Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err := docs.Parse(file); err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	parser := &Parser{
		Docs: docs,
		pr:   pr,
	}

	go parser.transform(pw)

	return parser, nil
}

func (p *Parser) Read(b []byte) (int, error) {
	return p.pr.Read(b)
}

func (p *Parser) Close() error {
	p.pr.Close()

	return p.Docs.Close()
}

func (p *Parser) transform(pw *io.PipeWriter) {
	zw := zip.NewWriter(pw)

	if err := p.Docs.SerializeEntries(zw); err != nil {
		zw.Close() // Free underlying resources only
		pw.CloseWithError(fmt.Errorf("lsif parser: Docs.SerializeEntries: %v", err))
		return
	}

	if err := zw.Close(); err != nil {
		pw.CloseWithError(fmt.Errorf("lsif parser: ZipWriter.Close: %v", err))
		return
	}

	pw.Close()
}
