// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package address

import "sort"

// GetCountryCodes returns all known country codes.
func GetCountryCodes() []string {
	return countryCodes
}

// CheckCountryCode checks whether a countryCode is valid.
func CheckCountryCode(countryCode string) bool {
	if len(countryCode) != 2 {
		return false
	}
	i := sort.SearchStrings(countryCodes, countryCode)
	if i >= len(countryCodes) {
		return false
	}

	return countryCodes[i] == countryCode
}
