package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPagination_ToOffsetLimit(t *testing.T) {
	tcs := map[string]struct {
		input               Input
		expOffset, expLimit int
	}{
		"page 0, size 0": {
			input: Input{
				Page: 0,
				Size: 0,
			},
			expOffset: 0,
			expLimit:  30,
		},
		"page 0, size 20": {
			input: Input{
				Page: 0,
				Size: 20,
			},
			expOffset: 0,
			expLimit:  20,
		},
		"page 1, size 50": {
			input: Input{
				Page: 1,
				Size: 50,
			},
			expOffset: 0,
			expLimit:  50,
		},
		"page 2, size 30": {
			input: Input{
				Page: 2,
				Size: 30,
			},
			expOffset: 30,
			expLimit:  30,
		},
		"page 3, size 40": {
			input: Input{
				Page: 3,
				Size: 40,
			},
			expOffset: 80,
			expLimit:  40,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given

			// When
			offset, limit := ToOffsetLimit(tc.input)

			// Then
			require.Equal(t, tc.expOffset, offset)
			require.Equal(t, tc.expLimit, limit)
		})
	}
}
