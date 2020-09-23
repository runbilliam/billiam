// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// Package address provides address validation and formatting.
package address

// Address represents a customer address.
type Address struct {
	Line1 string
	Line2 string
	// Sublocality is the neighborhood/suburb/district.
	Sublocality string
	// Locality is the city/village/post town.
	Locality string
	// Region is the state/province/prefecture.
	// An ISO code is used when available.
	Region string
	// PostalCode is the postal/zip/pin code.
	PostalCode string
	// CountryCode is the two-letter code as defined by CLDR.
	CountryCode string
}

// IsEmpty returns whether a is empty.
func (a Address) IsEmpty() bool {
	// An address must at least have a country code.
	return a.CountryCode == ""
}
