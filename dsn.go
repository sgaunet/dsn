// Package dsn provides utilities for parsing and manipulating database connection strings (DSN).
package dsn

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Error messages.
var (
	ErrInvalidDSNFormat  = errors.New("invalid DSN format")
	ErrInvalidParameters = errors.New("invalid parameters format, must be key=value")
)

// DSN represents a database source string.
type DSN struct {
	dsn        string
	scheme     string
	user       string
	password   string
	host       string
	port       string
	dbname     string
	parameters map[string]string
}

// New creates a new DSN from a string.
// Format: <scheme>://<user>:<password>@<host>:<port>/<dbname>?<parameters>.
//
//nolint:cyclop
func New(dsn string) (*DSN, error) {
	// Reject paths that start with / and don't have a scheme or host part
	if strings.HasPrefix(dsn, "/") && !strings.Contains(dsn, "://") {
		return nil, fmt.Errorf("%w", ErrInvalidDSNFormat)
	}

	// Support for simple hostname case
	if !strings.Contains(dsn, "://") && !strings.Contains(dsn, "@") && !strings.Contains(dsn, "/") {
		d := DSN{
			dsn:        dsn,
			host:       dsn,
			parameters: make(map[string]string),
		}
		return &d, nil
	}

	// Ensure scheme for url.Parse to work correctly
	parsableDSN := dsn
	if !strings.Contains(parsableDSN, "://") {
		parsableDSN = "placeholder://" + parsableDSN
	}

	u, err := url.Parse(parsableDSN)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDSNFormat, err.Error())
	}

	// Extract scheme, accounting for the placeholder
	scheme := u.Scheme
	if scheme == "placeholder" {
		scheme = ""
	}

	// Extract user and password
	user := ""
	password := ""
	if u.User != nil {
		user = u.User.Username()
		if pass, ok := u.User.Password(); ok {
			password = pass
		}
	}

	// Extract host and port
	host := u.Hostname()
	port := u.Port()

	// Extract path (dbname)
	dbname := strings.TrimPrefix(u.Path, "/")

	// Extract parameters
	parameters := make(map[string]string)
	for k, values := range u.Query() {
		if len(values) > 0 {
			parameters[k] = values[len(values)-1] // Use the last value if multiple exist
		}
	}

	d := DSN{
		dsn:        dsn,
		scheme:     scheme,
		user:       user,
		password:   password,
		host:       host,
		port:       port,
		dbname:     dbname,
		parameters: parameters,
	}
	return &d, nil
}

// GetUser returns the username component of the DSN.
func (d *DSN) GetUser() string {
	return d.user
}

// GetPassword returns the password component of the DSN.
func (d *DSN) GetPassword() string {
	return d.password
}

// GetHost returns the hostname component of the DSN.
func (d *DSN) GetHost() string {
	return d.host
}

// GetPort returns the port component of the DSN or the default port if not specified.
func (d *DSN) GetPort(defaultPort string) string {
	if d.port == "" {
		return defaultPort
	}
	return d.port
}

// GetPortInt returns the port component of the DSN as an integer or the default port if not specified or invalid.
func (d *DSN) GetPortInt(defaultPort int) int {
	portInt, err := strconv.Atoi(d.port)
	if err != nil {
		return defaultPort
	}
	return portInt
}

// GetParameter returns the value of a specific parameter or empty string if not found.
func (d *DSN) GetParameter(parameter string) string {
	if val, exists := d.parameters[parameter]; exists {
		return val
	}
	return ""
}

// GetParameters returns all parameters as a map.
func (d *DSN) GetParameters() map[string]string {
	paramsCopy := make(map[string]string, len(d.parameters))
	for k, v := range d.parameters {
		paramsCopy[k] = v
	}
	return paramsCopy
}

// GetDBName returns the database name component of the DSN.
func (d *DSN) GetDBName() string {
	return d.dbname
}

// GetPostgresURI returns a PostgreSQL connection string format.
func (d *DSN) GetPostgresURI() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.GetHost(),
		d.GetPort("5432"),
		d.GetUser(),
		d.GetPassword(),
		d.GetDBName(),
		d.GetParameter("sslmode"))

	// Add other parameters that need special handling
	if d.GetParameter("search_path") != "" {
		psqlInfo += " search_path=" + d.GetParameter("search_path")
	}

	return psqlInfo
}

// GetScheme returns the scheme component of the DSN.
func (d *DSN) GetScheme() string {
	return d.scheme
}

// String returns the original DSN string.
func (d *DSN) String() string {
	return d.dsn
}

// StringRedacted returns the DSN with the password redacted.
func (d *DSN) StringRedacted() string {
	// If there's no password, just return the original
	if d.password == "" {
		return d.dsn
	}

	// Otherwise, construct a redacted version
	var result strings.Builder

	if d.scheme != "" {
		result.WriteString(d.scheme)
		result.WriteString("://")
	}

	if d.user != "" {
		result.WriteString(d.user)
		result.WriteString(":****@")
	}

	result.WriteString(d.host)

	if d.port != "" {
		result.WriteString(":")
		result.WriteString(d.port)
	}

	if d.dbname != "" {
		result.WriteString("/")
		result.WriteString(d.dbname)
	}

	if len(d.parameters) > 0 {
		result.WriteString("?")

		paramStrs := make([]string, 0, len(d.parameters))
		for k, v := range d.parameters {
			paramStrs = append(paramStrs, k+"="+v)
		}

		result.WriteString(strings.Join(paramStrs, "&"))
	}

	return result.String()
}

// SetParameter sets or updates a parameter in the DSN.
//
func (d *DSN) SetParameter(key, value string) *DSN {
	d.parameters[key] = value
	return d
}
