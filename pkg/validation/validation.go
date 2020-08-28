// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package validation

const (
	// CodeRequired is used for missing values.
	CodeRequired = "required"
	// CodeInvalidValue is used for invalid values.
	CodeInvalidValue = "invalid_value"
	// CodeInvalidChoice is used for values not found in a predefined list.
	CodeInvalidChoice = "invalid_choice"
	// CodeNotUnique is used for values that are not unique (e.g. email in use).
	CodeNotUnique = "not_unique"
)

// Error represents a validation error.
type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// Required creates a required validation error.
func Required(message string) Error {
	return Error{CodeRequired, message}
}

// InvalidValue creates an invalid value validation error.
func InvalidValue(message string) Error {
	return Error{CodeInvalidValue, message}
}

// InvalidEmail creates an invalid choice error.
func InvalidChoice(message string) Error {
	return Error{CodeInvalidChoice, message}
}

// NotUnique creates a not unique error.
func NotUnique(message string) Error {
	return Error{CodeNotUnique, message}
}

// Errors maps a field path to a list of errors.
type Errors map[string][]error

// IsEmpty returns whether the list is empty.
func (e Errors) IsEmpty() bool {
	return e == nil || len(e) == 0
}

// Get gets the first error associated with the given path.
//
// If there are no associated errors, returns nil.
// To access multiple errors, use the map directly.
func (e Errors) Get(path string) error {
	if e == nil {
		return nil
	}
	es := e[path]
	if len(es) == 0 {
		return nil
	}
	return es[0]
}

// Set sets the path to the error. Replaces any existing errors.
func (e Errors) Set(path string, err error) {
	e[path] = []error{err}
}

// Add adds the error to the path. Apends to any existing errors.
func (e Errors) Add(path string, err error) {
	e[path] = append(e[path], err)
}

// Del deletes the values associated with path.
func (e Errors) Del(path string) {
	delete(e, path)
}

// Merge adds the errors optionally prefixed by the path.
//
// Use Merge to add a nested struct's error list
// to the parent struct's error list.
func (e Errors) Merge(path string, errors Errors) {
	for errpath, errs := range errors {
		if path != "" {
			errpath = path + "." + errpath
		}
		e[errpath] = errs
	}
}
