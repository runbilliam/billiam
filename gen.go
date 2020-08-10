// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// +build ignore

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("migrations"), vfsgen.Options{
		Filename:        "migrations.go",
		PackageName:     "billiam",
		BuildTags:       "!dev",
		VariableName:    "Migrations",
		VariableComment: "Migrations are database schema migrations, embedded by vfsgen.",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
