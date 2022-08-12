package main

import (
	"fmt"
	"testing"
)

func TestIsValidCmd(t *testing.T) {
	testcases := []struct {
		in       uint32
		expected bool
	}{
		{1, true},
		{55, false},
		{1000, false},
		{2147483653, true},
	}
	for _, tc := range testcases {

		t.Run(fmt.Sprintf("Got %v in %v", isValidCmd(tc.in), tc), func(t *testing.T) {
			if isValidCmd(tc.in) != tc.expected {
				t.Error("FFFAAAIILL")
			}
		})
	}
}
