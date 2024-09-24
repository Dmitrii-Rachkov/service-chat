package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"service-chat/internal/logger/prettylog"
)

func TestSetupLogger(t *testing.T) {
	// Опции для красивого логгера
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "nothing" {
				return slog.Attr{}
			}
			return a
		},
	}

	tests := []struct {
		name string
		env  string
		want *slog.Logger
	}{
		{
			name: "Local",
			env:  envLocal,
			want: slog.New(prettylog.NewHandler(opts)),
		},
		{
			name: "Dev",
			env:  envDev,
			want: slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		},
		{
			name: "Prod",
			env:  envProd,
			want: slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SetupLogger(tt.env)
			assert.IsType(t, tt.want, actual)
		})
	}
}
