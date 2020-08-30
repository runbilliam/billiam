// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// +build dev

package setup

import "net/http"

// Assets are setup assets, read from disk.
var Assets http.FileSystem = http.Dir("setup/assets")
