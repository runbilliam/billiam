// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// +build dev

package billiam

import "net/http"

// Migrations are database schema migrations, read from disk.
var Migrations http.FileSystem = http.Dir("migrations")
