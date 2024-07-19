package spinner_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/marceljaworski/cli-spinner/spinner"
	"github.com/stretchr/testify/assert"
)

func TestSpinnerStart(t *testing.T) {
	testCases := []struct {
		name     string
		duration time.Duration
		expects  string
	}{
		{
			name:     "should write correct values after 2 frames",
			duration: time.Millisecond * 25,
			expects:  "-\b\\\b",
		},
		{
			name:     "should write correct values after 6 frames",
			duration: time.Millisecond * 105,
			expects:  "-\b\\\b|\b/\b-\b\\\b",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			s := spinner.New(spinner.Config{
				Writer:    buf,
				FrameRate: time.Millisecond * 20,
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s.Start(ctx)
			time.Sleep(tc.duration)
			s.Stop()

			data, err := io.ReadAll(buf)
			assert.NoError(t, err)

			assert.Equal(t, tc.expects, string(data))
		})
	}
}

func TestSpinnerWorksAsync(t *testing.T) {
	s := spinner.New(spinner.Config{
		Writer:    &bytes.Buffer{},
		FrameRate: time.Millisecond * 5,
	})

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})

	go func() {
		s.Start(ctx)
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	select {
	case <-time.After(time.Millisecond * 200):
		assert.Fail(t, "test timed out")
	case <-done:
		//Test passed
	}
}

func TestWaitingAfterStop(t *testing.T) {
	buf := &bytes.Buffer{}

	s := spinner.New(spinner.Config{
		Writer:    buf,
		FrameRate: time.Millisecond * 20,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Start(ctx)
	time.Sleep(time.Millisecond * 35)
	s.Stop()
	time.Sleep(time.Millisecond * 10)

	data, err := io.ReadAll(buf)
	assert.NoError(t, err)

	assert.Equal(t, "-\b\\\b", string(data))
}

func TestSpinnerDoesNotPrintOnceStopped(t *testing.T) {
	s := spinner.New(spinner.Config{
		Writer:    &bytes.Buffer{},
		FrameRate: time.Millisecond * 5,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		s.Start(ctx)
		time.Sleep(10 * time.Millisecond)
		s.Stop()
		close(done)
	}()

	select {
	case <-time.After(time.Millisecond * 200):
		assert.Fail(t, "test timed out")
	case <-done:
		//Test passed
	}
}
