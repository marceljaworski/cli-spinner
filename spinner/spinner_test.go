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
			expects:  "-\b\\",
		},
		{
			name:     "should write correct values after 6 frames",
			duration: time.Millisecond * 105,
			expects:  "-\b\\\b|\b/\b-\b\\",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			s := spinner.New(spinner.Config{
				Writer:    buf,
				FrameRate: time.Millisecond * 20,
			})

			ctx, cancel := context.WithTimeout(context.Background(), tc.duration)
			defer cancel()

			s.Start(ctx)

			data, err := io.ReadAll(buf)
			assert.NoError(t, err)

			assert.Equal(t, tc.expects, string(data))
		})
	}
}
