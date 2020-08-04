// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

// Package logger provides zerolog helpers.
package logger

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/rs/zerolog"
)

// New creates a new zerolog logger.
//
// The logFormat is one of: text, json.
// The logLevel is one of: debug, info, warn, error, fatal.
func New(logFormat, logLevel string) (*zerolog.Logger, error) {
	var logger zerolog.Logger
	if logFormat == "text" {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: true,
		})
	} else {
		logger = logger.Output(os.Stderr)
	}
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("Unrecognized log level: %s", logLevel)
	}
	logger = logger.Level(level)
	// Add a timestamp to every log event.
	logger = logger.With().Timestamp().Logger()

	return &logger, nil
}

// NewStandard creates a new standard logger from the given zerolog logger.
//
// Used for packages which depend on the stdlib *log.Logger, e.g. net/http.
func NewStandard(logger *zerolog.Logger) *stdlog.Logger {
	return stdlog.New(logAdapter{logger}, "", 0)
}

type logAdapter struct {
	log *zerolog.Logger
}

// Write implements the io.Writer interface.
func (a logAdapter) Write(p []byte) (int, error) {
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[0 : n-1]
	}
	a.log.WithLevel(zerolog.NoLevel).Msg(string(p))

	return n, nil
}
