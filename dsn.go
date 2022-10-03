package dsn

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

type DSN struct {
	dsn string
}

func New(dsn string) *DSN {
	d := DSN{
		dsn: dsn,
	}
	return &d
}

func (d *DSN) GetUser() string {
	u, _ := url.Parse(d.dsn)
	return u.User.Username()
}

func (d *DSN) GetPassword() string {
	u, _ := url.Parse(d.dsn)
	password, _ := u.User.Password()
	return password
}

func (d *DSN) GetHost() string {
	u, _ := url.Parse(d.dsn)
	host, _, _ := net.SplitHostPort(u.Host)
	return host
}

func (d *DSN) GetPort() string {
	u, _ := url.Parse(d.dsn)
	_, port, _ := net.SplitHostPort(u.Host)
	return port
}

func (d *DSN) GetFragment(parameter string) string {
	var value string
	u, _ := url.Parse(d.dsn)
	m, _ := url.ParseQuery(u.RawQuery)
	if _, isKeyPresent := m[parameter]; isKeyPresent {
		value = m[parameter][0]
	}
	return value
}

func (d *DSN) GetPath() string {
	u, _ := url.Parse(d.dsn)
	return u.Path
}

func (d *DSN) GetPathWithoutSlash() string {
	u, _ := url.Parse(d.dsn)
	return strings.Replace(u.Path, "/", "", 1)
}

func (d *DSN) GetPostgresUri() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.GetHost(), d.GetPort(), d.GetUser(), d.GetPassword(), d.GetPathWithoutSlash(), d.GetFragment("sslmode"))
	if d.GetFragment("search_path") != "" {
		psqlInfo = psqlInfo + fmt.Sprintf(" search_path=%s", d.GetFragment("search_path"))
	}
	return psqlInfo
}
