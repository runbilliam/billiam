// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"testing"

	"github.com/runbilliam/billiam/pkg/validation"
)

func TestCheckEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"", false},
		{"example", false},
		{"@example.com", false},
		{"email@-example.com", false},
		{"email@example@example.com", false},
		{"email@example..com", false},
		{"Jane Smith <email@example.com>", false},
		{"#@%^%#$@#$@#.com", false},

		{"email@example.com", true},
		{"email@subdomain.example.com", true},
		{"email@example-japan.co.jp", true},
		{"email@example.city", true},
		{"fname.lname@example.com", true},
		{"fname.lname+throwaway@example.com", true},
		{"fname-lname@example.com", true},
		{"fname_lname@example.com", true},
		{"12345@example.com", true},
		{"_______@example.com", true},

		// Allowed by HTML5, disallowed by RFC 5322.
		{".email@example.com", true},
		{"email.@example.com", true},
		{"email..email@example.com", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := validation.CheckEmail(tt.email)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
