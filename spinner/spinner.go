package spinner

import (
	"context"
	"io"
	"os"
	"sync"
	"time"
)

type Spinner struct {
	writer     io.Writer
	frameRate  time.Duration
	frames     []rune
	cancelFunc context.CancelFunc
	doneCh     chan struct{}
	lock       sync.RWMutex
}

type Config struct {
	Writer    io.Writer
	FrameRate time.Duration
}

func New(cfg Config) *Spinner {
	s := &Spinner{
		writer:    os.Stderr,
		frameRate: time.Millisecond * 250,
		frames:    []rune{'-', '\\', '|', '/'},
	}
	if cfg.Writer != nil {
		s.writer = cfg.Writer
	}
	if cfg.FrameRate != 0 {
		s.frameRate = cfg.FrameRate
	}
	return s
}

func (s *Spinner) Start() {
	if s.isRunning() {
		return
	}

	s.lock.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	done := make(chan struct{})
	s.doneCh = done
	s.lock.Unlock()
	go func() {
		for {
			for _, frame := range s.frames {
				b := byte(frame)
				s.writer.Write([]byte{b})

				select {
				case <-ctx.Done():
					s.writer.Write([]byte("\b"))
					close(done)
					return
				case <-time.After(s.frameRate):
					break
				}

				s.writer.Write([]byte("\b"))
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if s.isRunning() {
		return
	}
	s.cancelFunc()
	<-s.doneCh

	s.lock.Lock()
	defer s.lock.Unlock()

	s.doneCh = nil
}

func (s *Spinner) isRunning() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.doneCh != nil
}
