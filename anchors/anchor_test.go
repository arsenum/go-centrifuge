// +build unit

package anchors

import (
	"testing"

	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewAnchorId(t *testing.T) {
	tests := []struct {
		name  string
		slice []byte
		err   string
	}{
		{
			"smallerSlice",
			utils.RandomSlice(AnchorIDLength - 1),
			"invalid length byte slice provided for anchorID",
		},
		{
			"largerSlice",
			utils.RandomSlice(AnchorIDLength + 1),
			"invalid length byte slice provided for anchorID",
		},
		{
			"nilSlice",
			nil,
			"invalid length byte slice provided for anchorID",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewAnchorID(test.slice)
			assert.Equal(t, test.err, err.Error())
		})
	}
}

func TestNewDocRoot(t *testing.T) {
	tests := []struct {
		name  string
		slice []byte
		err   string
	}{
		{
			"smallerSlice",
			utils.RandomSlice(RootLength - 1),
			"invalid length byte slice provided for docRoot",
		},
		{
			"largerSlice",
			utils.RandomSlice(RootLength + 1),
			"invalid length byte slice provided for docRoot",
		},
		{
			"nilSlice",
			nil,
			"invalid length byte slice provided for docRoot",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewDocRoot(test.slice)
			assert.Equal(t, test.err, err.Error())
		})
	}
}
