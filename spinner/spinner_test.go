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
	t.Run("should write correct values after 2 frames", func(t *testing.T) {
		buf := &bytes.Buffer{}

		s := spinner.New(spinner.Config{
			Writer:    buf,
			FrameRate: time.Millisecond * 20,
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*25)
		defer cancel()

		s.Start(ctx)

		data, err := io.ReadAll(buf)
		assert.NoError(t, err)

		assert.Equal(t, "-\b\\", string(data))
	})
}
