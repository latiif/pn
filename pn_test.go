package pn

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValid(t *testing.T) {
	tcs := []struct {
		description   string
		input         interface{}
		expectedError error
		expectedInfo  Info
	}{
		{
			description:   "Invalid input type",
			input:         true,
			expectedError: invalidInputValue(true),
		},
		{
			description:   "Invalid personal number (string)",
			input:         "950101",
			expectedError: invalidPn("950101"),
		},
		{
			description:   "Invalid personal number (non-string)",
			input:         950101,
			expectedError: invalidPn("950101"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			info, err := GetInfo(tc.input)

			if err != tc.expectedError {
				if diff := cmp.Diff(err.Error(), tc.expectedError.Error()); diff != "" {
					t.Logf("got: %v, want: %v", err, tc.expectedError)
					t.Fail()
				}
			}
			if !cmp.Equal(info, tc.expectedInfo) {
				t.Logf("diff: %v", cmp.Diff(info, tc.expectedInfo))
			}
		})
	}
}
