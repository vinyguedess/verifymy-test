package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceContains(t *testing.T) {
	tests := []struct {
		description    string
		slice          []interface{}
		element        interface{}
		expectedResult bool
	}{
		{
			description:    "Contains string",
			slice:          []interface{}{"hello", "world"},
			element:        "hello",
			expectedResult: true,
		},
		{
			description:    "Does not contain string",
			slice:          []interface{}{"hello", "world"},
			element:        "foo",
			expectedResult: false,
		},
		{
			description:    "Contains int",
			slice:          []interface{}{1, 2, 3},
			element:        2,
			expectedResult: true,
		},
		{
			description:    "Does not contain int",
			slice:          []interface{}{1, 2, 3},
			element:        4,
			expectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := SliceContains(test.slice, test.element)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}
