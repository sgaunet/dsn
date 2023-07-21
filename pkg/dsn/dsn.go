package dsn

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type DSN interface {
	GetUser() string
	GetPassword() string
	GetHost() string
	GetPort(string) string
	GetPortInt(int) int
	GetFragment(string) string
	GetPath() string
	GetPathWithoutSlash() string
	GetPostgresUri() string
	GetScheme() string
}

type dsntype struct {
	dsn string
}

func New(dsn string) (DSN, error) {
	r := regexp.MustCompile(`\w+://(\w+@|\w+:\w+@)?[^:/]*(:\d*|:\w+)?(/.*)?$`)
	if !r.MatchString(dsn) {
		return nil, errors.New("wrong format")
	}
	d := dsntype{
		dsn: dsn,
	}
	return &d, nil
}

func (d *dsntype) GetUser() string {
	u, _ := url.Parse(d.dsn)
	return u.User.Username()
}

func (d *dsntype) GetPassword() string {
	u, _ := url.Parse(d.dsn)
	password, _ := u.User.Password()
	return password
}

func (d *dsntype) GetHost() string {
	u, _ := url.Parse(d.dsn)
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return u.Host
	}
	return host
}

func (d *dsntype) GetPort(defaultPort string) string {
	u, _ := url.Parse(d.dsn)
	_, port, _ := net.SplitHostPort(u.Host)
	if port == "" {
		return defaultPort
	}
	return port
}

func (d *dsntype) GetPortInt(defaultPort int) int {
	r := regexp.MustCompile(`(\w|.)*:\d+.*`)
	if !r.MatchString(d.dsn) {
		return defaultPort
	}
	u, _ := url.Parse(d.dsn)
	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return defaultPort
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return defaultPort
	}
	return portInt
}

func (d *dsntype) GetFragment(parameter string) string {
	var value string
	u, _ := url.Parse(d.dsn)
	m, _ := url.ParseQuery(u.RawQuery)
	if _, isKeyPresent := m[parameter]; isKeyPresent {
		value = m[parameter][0]
	}
	return value
}

func (d *dsntype) GetPath() string {
	u, _ := url.Parse(d.dsn)
	return u.Path
}

func (d *dsntype) GetPathWithoutSlash() string {
	u, _ := url.Parse(d.dsn)
	return strings.Replace(u.Path, "/", "", 1)
}

func (d *dsntype) GetPostgresUri() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.GetHost(), d.GetPort("5432"), d.GetUser(), d.GetPassword(), d.GetPathWithoutSlash(), d.GetFragment("sslmode"))
	if d.GetFragment("search_path") != "" {
		psqlInfo = psqlInfo + fmt.Sprintf(" search_path=%s", d.GetFragment("search_path"))
	}
	return psqlInfo
}

func (d *dsntype) GetScheme() string {
	u, _ := url.Parse(d.dsn)
	return u.Scheme
}
