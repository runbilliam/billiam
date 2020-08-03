// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: Apache-2.0

package billiam

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bojanz/envx"
	"github.com/pelletier/go-toml"
)

const exampleConfig = `
# Each quoted value can contain one or more environment variables.
# For example, a site deployed on Heroku would use "${PORT}" and "${DATABASE_URL}".
# Each environment variable can have a default value:
#   ${LISTEN:2490} defaults to 2490 if $LISTEN is not set.
[server]
listen = "${LISTEN:2490}" # port number or systemd socket name.
tls_listen = "${TLS_LISTEN:2491}" # port number or systemd socket name.
tls_cert = "${TLS_CERT}" # path to a cert.pem
tls_key = "${TLS_KEY}" # path to a key.pem

[database]
url = "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT:5432}/${DB_NAME}"

[log]
format = "${LOG_FORMAT:json}" # One of: text, json.
level = "${LOG_LEVEL:info}" # One of: debug, info, warn, error, fatal.
`

// Config represents the app configuration.
type Config struct {
	Server struct {
		Listen    string
		TLSListen string `toml:"tls_listen"`
		TLSCert   string `toml:"tls_cert"`
		TLSKey    string `toml:"tls_key"`
	}
	Database struct {
		URL string
	}
	Log struct {
		Format string
		Level  string
	}
}

// CreateConfig creates a config file with the given filename.
func CreateConfig(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.WriteString(exampleConfig)
	if err != nil {
		return err
	}

	return f.Close()
}

// ReadConfig reads a config file with the given filename.
func ReadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%s is empty", filename)
	}
	config := &Config{}
	err = toml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	// Each value is allowed to contain environment variables.
	config.Server.Listen = envx.Expand(config.Server.Listen)
	config.Server.TLSListen = envx.Expand(config.Server.TLSListen)
	config.Server.TLSCert = envx.Expand(config.Server.TLSCert)
	config.Server.TLSKey = envx.Expand(config.Server.TLSKey)
	config.Database.URL = envx.Expand(config.Database.URL)
	config.Log.Format = envx.Expand(config.Log.Format)
	config.Log.Level = envx.Expand(config.Log.Level)

	return config, nil
}
