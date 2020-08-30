// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/bojanz/currency"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/vfsutil"

	"github.com/runbilliam/billiam/internal/settings"
	"github.com/runbilliam/billiam/internal/user"
	"github.com/runbilliam/billiam/pkg/timezone"
	"github.com/runbilliam/billiam/pkg/validation"
)

type data struct {
	Values           url.Values
	Errors           validation.Errors
	Timezones        []string
	CommonCurrencies []string
	OtherCurrencies  []string
}

// Handler handles setup routes.
type Handler struct {
	logger *zerolog.Logger
}

// NewHandler creates a new setup handler.
func NewHandler(logger *zerolog.Logger) *Handler {
	h := Handler{
		logger: logger,
	}
	return &h
}

// Routes attaches setup routes to the router.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.Setup)
	r.Post("/", h.SubmitSetup)
	r.Get("/assets/*", h.ServeAssets)
}

// Setup renders the setup page.
func (h *Handler) Setup(w http.ResponseWriter, r *http.Request) {
	currencies := currency.GetCurrencyCodes()
	data := data{
		Timezones:        timezone.GetNames(),
		CommonCurrencies: currencies[0:10],
		OtherCurrencies:  currencies[10:],
	}
	h.render(w, "setup.html.hbs", data)
}

// SubmitSetup handles the setup page submit.
func (h *Handler) SubmitSetup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.handleError(w, err)
		return
	}
	u := user.New()
	u.Email = strings.TrimSpace(r.PostFormValue("email"))
	u.Password = r.PostFormValue("password")
	u.Timezone = r.PostFormValue("timezone")
	u.Active = true
	errs := u.Validate()

	s := settings.New()
	s.Timezone = r.PostFormValue("timezone")
	s.Currency = r.PostFormValue("currency")
	errs.Merge("", s.Validate())

	if !errs.IsEmpty() {
		currencies := currency.GetCurrencyCodes()
		data := data{
			Values:           r.PostForm,
			Errors:           errs,
			Timezones:        timezone.GetNames(),
			CommonCurrencies: currencies[0:10],
			OtherCurrencies:  currencies[10:],
		}
		h.render(w, "setup.html.hbs", data)
		return
	}

	w.Write([]byte("Setup complete!"))
}

// ServeAssets serves assets.
func (h *Handler) ServeAssets(w http.ResponseWriter, r *http.Request) {
	// Skip directory listings and templates.
	fs := filter.Skip(Assets, func(path string, fi os.FileInfo) bool {
		return fi.IsDir() || strings.HasPrefix("templates/", path)
	})
	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fsHandler := http.StripPrefix(pathPrefix, http.FileServer(fs))
	fsHandler.ServeHTTP(w, r)
}

func (h *Handler) render(w http.ResponseWriter, filename string, data interface{}) {
	b, err := vfsutil.ReadFile(Assets, "templates/"+filename)
	if err != nil {
		h.handleError(w, err)
	}
	tpl, err := raymond.Parse(string(b))
	if err != nil {
		err := fmt.Errorf("%v: %w", filename, err)
		h.handleError(w, err)
		return
	}
	result, err := tpl.Exec(data)
	if err != nil {
		err := fmt.Errorf("%v: %w", filename, err)
		h.handleError(w, err)
		return
	}
	w.Write([]byte(result))
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	h.logger.Error().Msg(err.Error())
	http.Error(w, "Internal Server Error", 500)
}
