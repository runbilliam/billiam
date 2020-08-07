// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package billiam

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bojanz/httpx"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"

	"github.com/runbilliam/billiam/pkg/log"
)

// Version is the current application version. Replaced at build time.
var Version = "v1"

// Application represents the application.
type Application struct {
	cfg    *Config
	logger *zerolog.Logger
	db     *pgxpool.Pool
	server *httpx.Server
}

// New creates a new application.
func New(cfg *Config, logger *zerolog.Logger, db *pgxpool.Pool) (*Application, error) {
	var server *httpx.Server
	stdLogger := log.NewStandard(logger)
	if cfg.Server.TLSCert != "" {
		httpsAddr := toAddr(cfg.Server.TLSListen)
		cert, err := tls.LoadX509KeyPair(cfg.Server.TLSCert, cfg.Server.TLSKey)
		if err != nil {
			return nil, err
		}
		server = httpx.NewServerTLS(httpsAddr, cert, nil)
		server.ErrorLog = stdLogger
	} else {
		httpAddr := toAddr(cfg.Server.Listen)
		server = httpx.NewServer(httpAddr, nil)
		server.ErrorLog = stdLogger
	}
	app := &Application{
		cfg:    cfg,
		logger: logger,
		db:     db,
		server: server,
	}

	return app, nil
}

// Start starts the application.
func (app *Application) Start() error {
	app.logger.Info().Msgf("Starting billiam %s", Version)
	app.server.Handler = app.buildRouter()

	proto := "HTTP"
	if app.server.IsTLS() {
		proto = "HTTPS"
	}
	app.logger.Info().Msgf("Listening for %v on %v", proto, app.server.Addr)
	if err := app.server.Start(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown shuts down the application.
func (app *Application) Shutdown() error {
	app.logger.Info().Msgf("Shutting down")

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := app.server.Shutdown(ctx)
	if err == context.DeadlineExceeded {
		proto := "HTTP"
		if app.server.IsTLS() {
			proto = "HTTPS"
		}
		return fmt.Errorf("%v timeout exceeded while waiting on %v shutdown", timeout, proto)
	}

	return nil
}

// buildRouter builds the router.
func (app *Application) buildRouter() *chi.Mux {
	if app.cfg.Log.Format == "json" {
		httplog.DefaultOptions.JSON = true
	}
	// Log only responses (default is request&response).
	httplog.DefaultOptions.Concise = true

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(*app.logger))
	r.Use(middleware.Heartbeat("/health"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return r
}

// toAddr() converts a port number / systemd socket name into an addr.
func toAddr(listen string) string {
	if listen == "" {
		return ""
	}

	var addr string
	if _, err := strconv.Atoi(listen); err == nil {
		// Port number. Prefix with ":".
		addr = ":" + listen
	} else {
		// Systemd socket. Prefix with "systemd:".
		addr = "systemd:" + listen
	}

	return addr
}
