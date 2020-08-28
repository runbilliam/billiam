// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import "regexp"

// https://html.spec.whatwg.org/multipage/forms.html#valid-e-mail-address
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// CheckEmail checks whether the given email is valid.
func CheckEmail(email string) bool {
	if email == "" {
		return false
	}
	return len(email) <= 254 && rxEmail.MatchString(email)
}
