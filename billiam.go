// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package billiam

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/bojanz/httpx"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/runbilliam/billiam/pkg/log"
)

// Version is the current application version. Replaced at build time.
var Version = "v1"

// Application represents the application.
type Application struct {
	cfg            *Config
	logger         *zerolog.Logger
	db             *pgxpool.Pool
	mainServer     *httpx.Server
	redirectServer *httpx.Server
}

// New creates a new application.
func New(cfg *Config, logger *zerolog.Logger, db *pgxpool.Pool) (*Application, error) {
	// Initialize the HTTP servers.
	var mainServer, redirectServer *httpx.Server
	httpAddr := toAddr(cfg.Server.Listen)
	httpsAddr := toAddr(cfg.Server.TLSListen)
	stdLogger := log.NewStandard(logger)
	if cfg.Server.TLSCert != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Server.TLSCert, cfg.Server.TLSKey)
		if err != nil {
			return nil, err
		}
		mainServer = httpx.NewServerTLS(httpsAddr, cert, nil)
		mainServer.ErrorLog = stdLogger

		redirectServer = httpx.NewServer(httpAddr, httpRedirectHandler{})
		redirectServer.ErrorLog = stdLogger
	} else {
		mainServer = httpx.NewServer(httpAddr, nil)
		mainServer.ErrorLog = stdLogger
	}
	app := &Application{
		cfg:            cfg,
		logger:         logger,
		db:             db,
		mainServer:     mainServer,
		redirectServer: redirectServer,
	}

	return app, nil
}

// Start starts the application.
func (app *Application) Start() error {
	app.logger.Info().Msgf("Starting billiam %s", Version)
	app.mainServer.Handler = app.buildRouter()

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		proto := "HTTP"
		if app.mainServer.IsTLS() {
			proto = "HTTPS"
		}
		app.logger.Info().Msgf("Listening for %v on %v", proto, app.mainServer.Addr)
		if err := app.mainServer.Start(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	if app.redirectServer != nil {
		g.Go(func() error {
			app.logger.Info().Msgf("Listening for HTTP on %v", app.redirectServer.Addr)
			if err := app.redirectServer.Start(); err != http.ErrServerClosed {
				return err
			}
			return nil
		})
	}
	go func() {
		// The context is closed if both servers finish, or one of them
		// errors out, in which case we want to close the other and return.
		<-ctx.Done()
		app.mainServer.Close()
		if app.redirectServer != nil {
			app.redirectServer.Close()
		}
	}()

	return g.Wait()
}

// Shutdown shuts down the application.
func (app *Application) Shutdown() error {
	app.logger.Info().Msgf("Shutting down")

	if app.redirectServer != nil {
		redirectTimeout := 1 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), redirectTimeout)
		defer cancel()
		if err := app.redirectServer.Shutdown(ctx); err == context.DeadlineExceeded {
			return fmt.Errorf("%v timeout exceeded while waiting on HTTP shutdown", redirectTimeout)
		}
	}
	mainTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), mainTimeout)
	defer cancel()
	if err := app.mainServer.Shutdown(ctx); err == context.DeadlineExceeded {
		proto := "HTTP"
		if app.mainServer.IsTLS() {
			proto = "HTTPS"
		}
		return fmt.Errorf("%v timeout exceeded while waiting on %v shutdown", mainTimeout, proto)
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

// httpRedirectHandler sends all HTTP traffic to the HTTPS server.
type httpRedirectHandler struct{}

// ServeHTTP implements the http.Handler interface.
func (h httpRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// An HTTPS port is not specified to avoid problems with Docker port forwarding.
	// Thus, the redirect will only work if Billiam is listening on standard ports.
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		// No port found.
		host = r.Host
	}
	r.URL.Host = host
	r.URL.Scheme = "https"

	w.Header().Set("Connection", "close")
	http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
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
