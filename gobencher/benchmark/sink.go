package benchmark

import (
	"io"
	"sync"
)

type Sink struct {
	mu  sync.Mutex
	sum int64
	tee io.Writer
}

func NewSink() *Sink {
	return &Sink{}
}

func (s *Sink) Write(p []byte) (int, error) {
	s.mu.Lock()
	s.sum += int64(len(p))
	tee := s.tee
	s.mu.Unlock()

	if tee != nil {
		if _, err := tee.Write(p); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

func (s *Sink) Sync() error {
	return nil
}

func (s *Sink) Reset() {
	s.mu.Lock()
	s.sum = 0
	s.mu.Unlock()
}

func (s *Sink) BytesWritten() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sum
}

func (s *Sink) SetTee(w io.Writer) {
	s.mu.Lock()
	s.tee = w
	s.mu.Unlock()
}
