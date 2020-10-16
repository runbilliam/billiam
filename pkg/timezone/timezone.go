// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// Package timezone provides a list of timezones.
package timezone

import (
	"sort"
)

// GetNames returns a list of timezone names.
func GetNames() []string {
	return names
}

// IsValid checks whether a timezone name is valid.
//
// An empty timezone name is considered valid.
func IsValid(name string) bool {
	if name == "" {
		return true
	}
	i := sort.SearchStrings(names, name)
	if i >= len(names) {
		return false
	}

	return names[i] == name
}
