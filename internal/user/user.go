// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/runbilliam/billiam/pkg/timezone"
	"github.com/runbilliam/billiam/pkg/validation"
)

// ErrNotFound is returned when a user could not be found.
var ErrNotFound = errors.New("user not found")

type User struct {
	ID        ulid.ULID `json:"id"`
	Version   int       `json:"version"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Timezone  string    `json:"timezone"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LoginAt   time.Time `json:"login_at"`
}

// New creates a new user.
func New() User {
	now := time.Now().UTC()
	u := User{
		ID:        ulid.MustNew(ulid.Timestamp(now), rand.Reader),
		Version:   1,
		CreatedAt: now,
	}

	return u
}

// Validate validates the user.
func (u User) Validate() validation.Errors {
	errs := validation.Errors{}
	if u.ID == (ulid.ULID{}) {
		errs.Add("id", validation.Required("ID is required."))
	}
	if u.Version == 0 {
		errs.Add("version", validation.Required("Version is required."))
	}
	if u.Email == "" {
		errs.Add("email", validation.Required("Email is required."))
	}
	if u.Password == "" {
		errs.Add("password", validation.Required("Password is required."))
	}
	if u.Timezone == "" {
		errs.Add("timezone", validation.Required("Timezone is required."))
	}
	if u.CreatedAt.IsZero() {
		errs.Add("created_at", validation.Required("CreatedAt is required."))
	}

	if !validation.CheckEmail(u.Email) {
		errs.Add("email", validation.InvalidValue("Email is invalid."))
	}
	if !timezone.IsValid(u.Timezone) {
		errs.Add("timezone", validation.InvalidChoice("Invalid timezone."))
	}

	return errs
}
