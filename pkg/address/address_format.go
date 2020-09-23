// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package address

import "regexp"

// Field represents the address field.
type Field string

const (
	FieldLine1       Field = "1"
	FieldLine2       Field = "2"
	FieldSublocality Field = "S"
	FieldLocality    Field = "L"
	FieldRegion      Field = "R"
	FieldPostalCode  Field = "P"
)

// SublocalityType represents the sublocality type.
type SublocalityType uint8

const (
	SublocalityTypeSuburb SublocalityType = iota
	SublocalityTypeDistrict
	SublocalityTypeNeighborhood
	SublocalityTypeVillageTownship
	SublocalityTypeTownland
)

// LocalityType represents the locality type.
type LocalityType uint8

const (
	LocalityTypeCity LocalityType = iota
	LocalityTypeDistrict
	LocalityTypePostTown
	LocalityTypeSuburb
)

// RegionType represents the region type.
type RegionType uint8

const (
	RegionTypeProvince RegionType = iota
	RegionTypeArea
	RegionTypeCanton
	RegionTypeCounty
	RegionTypeDepartment
	RegionTypeDistrict
	RegionTypeDoSi
	RegionTypeEmirate
	RegionTypeIsland
	RegionTypeOblast
	RegionTypeParish
	RegionTypePrefecture
	RegionTypeState
)

// PostalCodeType represents the postal code type.
type PostalCodeType uint8

const (
	PostalCodeTypePostal PostalCodeType = iota
	PostalCodeTypeEir
	PostalCodeTypePin
	PostalCodeTypeZip
)

// Format represents an address format.
type Format struct {
	Layout            string
	Required          []Field
	SublocalityType   SublocalityType
	LocalityType      LocalityType
	RegionType        RegionType
	PostalCodeType    PostalCodeType
	PostalCodePattern string
	ShowRegionID      bool
	Regions           map[string]string
}

// Requires returns whether the given field is required.
func (f Format) IsRequired(field Field) bool {
	for _, ff := range f.Required {
		if ff == field {
			return true
		}
	}
	return false
}

// CheckRegion checks whether the given region is valid.
func (f Format) CheckRegion(region string) bool {
	if region == "" {
		return false
	}
	if len(f.Regions) == 0 {
		// A non-empty region is always valid if no regions have been predefined.
		return true
	}
	_, ok := f.Regions[region]
	return ok
}

// CheckPostalCode checks whether the given postal code is valid.
func (f Format) CheckPostalCode(postalCode string) bool {
	if postalCode == "" {
		return false
	}
	if f.PostalCodePattern == "" {
		// A non-empty postal code is always valid if no pattern is available.
		return true
	}
	rx := regexp.MustCompile(f.PostalCodePattern)
	return rx.MatchString(postalCode)
}

// GetFormat returns an address format for the given country code.
func GetFormat(countryCode string) Format {
	format, ok := formats[countryCode]
	if !ok {
		return formats["ZZ"]
	}
	return format
}
