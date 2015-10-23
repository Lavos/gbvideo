package gbvideo

import (
	"io"
)

type ProgressReader struct {
	reader io.Reader
	BytesRead chan int
}

func NewProgressReader (reader io.Reader) *ProgressReader {
	return &ProgressReader{
		reader: reader,
		BytesRead: make(chan int),
	}
}

func (pr *ProgressReader) Read (p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.BytesRead <- n

	if err != nil {
		close(pr.BytesRead)
	}

	return n, err
}
