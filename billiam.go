// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package billiam

import (
	"context"
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

	"github.com/runbilliam/billiam/pkg/logger"
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
func New(cfg *Config, logger *zerolog.Logger) (*Application, error) {
	db, err := pgxpool.Connect(context.Background(), cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	app := &Application{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	return app, nil
}

// Start starts the application.
func (app *Application) Start() error {
	app.logger.Info().Msgf("Starting billiam %s", Version)
	httpAddr := toAddr(app.cfg.Server.Listen)
	httpsAddr := toAddr(app.cfg.Server.TLSListen)
	stdLogger := logger.NewStandard(app.logger)
	r := app.buildRouter()

	if app.cfg.Server.TLSCert != "" {
		app.logger.Info().Msgf("Listening for HTTPS on %v", httpsAddr)
		app.server = httpx.NewServer(httpsAddr, r)
		app.server.ErrorLog = stdLogger
		err := app.server.ListenAndServeTLS(app.cfg.Server.TLSCert, app.cfg.Server.TLSKey)
		if err != http.ErrServerClosed {
			return err
		}
	} else {
		app.logger.Info().Msgf("Listening for HTTP on %v", httpAddr)
		app.server = httpx.NewServer(httpAddr, r)
		app.server.ErrorLog = stdLogger
		err := app.server.ListenAndServe()
		if err != http.ErrServerClosed {
			return err
		}
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
		return fmt.Errorf("%v timeout exceeded while waiting on shutdown", timeout)
	}
	app.db.Close()

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
