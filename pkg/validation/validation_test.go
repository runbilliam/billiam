// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"testing"

	"github.com/runbilliam/billiam/pkg/validation"
)

func TestErrors_IsEmpty(t *testing.T) {
	// Nil map.
	var errs validation.Errors
	if errs.IsEmpty() != true {
		t.Error("got false, want true")
	}

	// Empty map.
	errs = validation.Errors{}
	if errs.IsEmpty() != true {
		t.Error("got false, want true")
	}

	errs.Add("email", validation.NotUnique("Email is already in use."))
	if errs.IsEmpty() != false {
		t.Error("got true, want false")
	}

	errs.Del("email")
	if errs.IsEmpty() != true {
		t.Error("got false, want true")
	}
}

func TestErrors_Get(t *testing.T) {
	// Nil map.
	var errs validation.Errors
	err := errs.Get("country")
	if err != nil {
		t.Errorf("got %+v, want nil", err)
	}

	// Empty map.
	errs = validation.Errors{}
	err = errs.Get("country")
	if err != nil {
		t.Errorf("got %+v, want nil", err)
	}

	// Multiple errors on the same key.
	errRequired := validation.Required("Country is required.")
	errInvalidChoice := validation.InvalidChoice("Invalid country.")
	errs.Add("country", errRequired)
	errs.Add("country", errInvalidChoice)
	err = errs.Get("country")
	if err != errRequired {
		t.Errorf("got %#v, want %#v", err, errRequired)
	}
}

func TestErrors_Set(t *testing.T) {
	errRequired := validation.Required("Country is required.")
	errInvalidValue := validation.InvalidValue("Invalid country.")

	// Empty map.
	errs := validation.Errors{}
	errs.Set("country", errRequired)
	err := errs.Get("country")
	if err != errRequired {
		t.Errorf("got %#v, want %#v", err, errRequired)
	}

	// Existing map.
	errs = validation.Errors{}
	errs.Add("country", errRequired)
	err = errs.Get("country")
	if err != errRequired {
		t.Errorf("got %#v, want %#v", err, errRequired)
	}
	errs.Set("country", errInvalidValue)
	err = errs.Get("country")
	if err != errInvalidValue {
		t.Errorf("got %#v, want %#v", err, errInvalidValue)
	}
}

func TestErrors_Merge(t *testing.T) {
	errRequired := validation.Required("Country is required.")
	errInvalidChoice := validation.InvalidChoice("Invalid country.")
	errNotUnique := validation.NotUnique("Email is already in use.")
	addrErrs := validation.Errors{}
	addrErrs.Add("country", errRequired)
	addrErrs.Add("country", errInvalidChoice)
	emailErrs := validation.Errors{}
	emailErrs.Add("email", errNotUnique)

	// Merging under a key.
	errs := validation.Errors{}
	errs.Merge("address", addrErrs)
	err := errs.Get("address.country")
	if err != errRequired {
		t.Errorf("got %#v, want %#v", err, errRequired)
	}
	n := len(errs["address.country"])
	if n != 2 {
		t.Errorf("got %#v, want 2", n)
	}

	// Merging without a key.
	errs.Merge("", emailErrs)
	err = errs.Get("email")
	if err != errNotUnique {
		t.Errorf("got %#v, want %#v", err, errNotUnique)
	}
}
