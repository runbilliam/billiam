// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"github.com/bojanz/currency"

	"github.com/runbilliam/billiam/pkg/timezone"
	"github.com/runbilliam/billiam/pkg/validation"
)

// Settings represent the site configuration.
type Settings struct {
	Version  int    `json:"version"`
	SiteName string `json:"site_name"`
	Timezone string `json:"timezone"`
	Currency string `json:"currency"`
}

// New creates new settings.
func New() Settings {
	s := Settings{
		Version:  1,
		SiteName: "Billiam",
		Timezone: "Europe/Berlin",
		Currency: "EUR",
	}

	return s
}

// Validate validates the settings.
func (s Settings) Validate() validation.Errors {
	errs := validation.Errors{}
	if s.Version == 0 {
		errs.Add("version", validation.Required("Version is required."))
	}
	if s.SiteName == "" {
		errs.Add("site_name", validation.Required("Site name is required."))
	}
	if s.Timezone == "" {
		errs.Add("timezone", validation.Required("Timezone is required."))
	}
	if s.Currency == "" {
		errs.Add("currency", validation.Required("Currency is required."))
	}

	if !timezone.IsValid(s.Timezone) {
		errs.Add("timezone", validation.InvalidChoice("Invalid timezone."))
	}
	if !currency.IsValid(s.Currency) {
		errs.Add("currency", validation.InvalidChoice("Invalid currency."))
	}

	return errs
}
