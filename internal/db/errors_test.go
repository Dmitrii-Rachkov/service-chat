package db

import (
	"database/sql"
	errDriver "database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pureErr(t *testing.T) {
	tests := []struct {
		name        string
		errIn       error
		errExpected error
	}{
		{
			name:        "Nil error",
			errIn:       nil,
			errExpected: nil,
		},
		{
			name:        "Error with pq",
			errIn:       errors.New("pq: Not found messages"),
			errExpected: errors.New(fmt.Sprintf("Not found messages")),
		},
		{
			name:        "Error with sql",
			errIn:       sql.ErrNoRows,
			errExpected: errors.New("no rows in result set"),
		},
		{
			name:        "Error with driver",
			errIn:       errDriver.ErrBadConn,
			errExpected: errors.New("bad connection"),
		},
		{
			name:        "Error without replace",
			errIn:       errors.New("no replacement"),
			errExpected: errors.New("no replacement"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pureErr(tt.errIn)
			assert.Equal(t, tt.errExpected, err)
		})
	}
}
