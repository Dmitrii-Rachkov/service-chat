package logger

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want slog.Attr
	}{
		{
			name: "OK",
			err:  errors.New("test error"),
			want: slog.Attr{
				Key:   "error",
				Value: slog.StringValue("test error"),
			},
		},
		{
			name: "Empty error",
			err:  errors.New(""),
			want: slog.Attr{
				Key:   "error",
				Value: slog.StringValue(""),
			},
		},
		{
			name: "Nil error",
			want: slog.Attr{
				Key:   "error",
				Value: slog.StringValue(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Err(tt.err))
		})
	}
}
