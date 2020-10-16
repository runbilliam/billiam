// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package timezone_test

import (
	"testing"

	"github.com/runbilliam/billiam/pkg/timezone"
)

func TestGetNames(t *testing.T) {
	names := timezone.GetNames()
	var got [3]string
	copy(got[:], names[0:3])
	want := [3]string{"Africa/Abidjan", "Africa/Accra", "Africa/Algiers"}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", true},
		{"INVALID", false},
		{"europe/belgrade", false},
		{"Europe/Belgrade", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := timezone.IsValid(tt.name)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
