// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// Package address provides address validation and formatting.
package address

import (
	"github.com/runbilliam/billiam/pkg/validation"
)

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

// Validate validates the address.
func (a Address) Validate() validation.Errors {
	errs := validation.Errors{}
	if a.CountryCode == "" {
		errs.Add("country_code", validation.Required("Country is required."))
		return errs
	} else if !CheckCountryCode(a.CountryCode) {
		errs.Add("country_code", validation.InvalidChoice("Invalid country."))
		return errs
	}

	format := GetFormat(a.CountryCode)
	if a.Line1 == "" && format.IsRequired(FieldLine1) {
		errs.Add("line1", validation.Required("Line1 is required."))
	}
	if a.Line2 == "" && format.IsRequired(FieldLine1) {
		errs.Add("line2", validation.Required("Line2 is required."))
	}
	if a.Sublocality == "" && format.IsRequired(FieldSublocality) {
		errs.Add("sublocality", validation.Required("Sublocality is required."))
	}
	if a.Locality == "" && format.IsRequired(FieldLocality) {
		errs.Add("locality", validation.Required("Locality is required."))
	}
	if a.Region == "" && format.IsRequired(FieldRegion) {
		errs.Add("region", validation.Required("Region is required."))
	} else if a.Region != "" && !format.CheckRegion(a.Region) {
		errs.Add("region", validation.InvalidChoice("Invalid region."))
	}

	if a.PostalCode == "" && format.IsRequired(FieldPostalCode) {
		errs.Add("postal_code", validation.Required("Postal code is required."))
	} else if a.PostalCode != "" && !format.CheckPostalCode(a.PostalCode) {
		errs.Add("postal_code", validation.InvalidChoice("Invalid postal code."))
	}

	return errs
}
