package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_ShutDown(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want error
	}{
		{
			name: "Success",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "Nil context",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				httpServer: &http.Server{
					Addr: ":8080",
				},
			}
			assert.Equal(t, tt.want, srv.ShutDown(tt.ctx))
		})
	}
}
