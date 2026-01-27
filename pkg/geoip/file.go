package geoip

import (
	"bytes"
	"embed"
	"io"
	"sync"
)

type EmbeddedFile struct {
	data   []byte
	reader *bytes.Reader
	mu     sync.RWMutex
	closed bool
}

func NewEmbeddedFile(fs embed.FS, path string) (*EmbeddedFile, error) {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &EmbeddedFile{
		data:   data,
		reader: bytes.NewReader(data),
	}, nil
}

func (f *EmbeddedFile) Read(p []byte) (n int, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.closed {
		return 0, io.ErrClosedPipe
	}

	return f.reader.Read(p)
}

func (f *EmbeddedFile) ReadAt(p []byte, off int64) (n int, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.closed {
		return 0, io.ErrClosedPipe
	}

	return f.reader.ReadAt(p, off)
}

func (f *EmbeddedFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return io.ErrClosedPipe
	}

	f.closed = true
	f.data = nil
	f.reader = nil
	return nil
}

func (f *EmbeddedFile) Seek(offset int64, whence int) (int64, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.closed {
		return 0, io.ErrClosedPipe
	}

	return f.reader.Seek(offset, whence)
}
