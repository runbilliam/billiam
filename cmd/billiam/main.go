// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/runbilliam/billiam"
	"github.com/runbilliam/billiam/pkg/log"
)

const usage = `
Usage: billiam [command]

Commands:
  init         Initialize a new site in the current directory
  serve        Start the HTTP server
  version      Show version information
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stdout, "Billiam is a headless recuring billing system")
		fmt.Fprintln(os.Stdout, usage)
		return
	}
	cmd := os.Args[1]

	switch cmd {
	case "init":
		cmdInit()
	case "serve":
		cmdServe()
	case "version":
		cmdVersion()
	default:
		fmt.Fprintln(os.Stderr, "Error: Unknown command", cmd)
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
}

func cmdInit() {
	if _, err := os.Stat("config.toml"); err == nil {
		fmt.Fprintln(os.Stderr, "Error: Site already initialized")
		os.Exit(1)
	}
	if err := billiam.CreateConfig("config.toml"); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	dir, _ := os.Getwd()
	fmt.Fprintln(os.Stdout, "Initialized a new site in", dir)
}

func cmdServe() {
	config, err := billiam.ReadConfig("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Error: No site found in the current directory")
		} else {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
		os.Exit(1)
	}
	logger, err := log.New(config.Log.Format, config.Log.Level)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	db, err := pgxpool.Connect(context.Background(), config.Database.URL)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	app, err := billiam.New(config, logger, db)
	if err != nil {
		db.Close()
		logger.Fatal().Msg(err.Error())
	}

	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal)
		signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
		<-shutdownCh

		errCh <- app.Shutdown()
	}()

	if err := app.Start(); err != nil {
		db.Close()
		logger.Fatal().Msg(err.Error())
	}
	if err := <-errCh; err != nil {
		logger.Warn().Msg(err.Error())
	}
	db.Close()
}

func cmdVersion() {
	fmt.Fprintf(os.Stdout, "billiam %s %s/%s %s\n",
		billiam.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
