package dsn

import (
	"errors"
	"fmt"
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
	GetParameter(string) string
	GetDBName() string
	GetPostgresUri() string
	GetScheme() string
	String() string
}

type dsntype struct {
	dsn        string
	scheme     string
	user       string
	password   string
	host       string
	port       string
	dbname     string
	parameters map[string]string
}

func New(dsn string) (DSN, error) {
	r := regexp.MustCompile(`^((?P<scheme>\w+)://)?((?P<user>[\w-]*):(?P<password>.*)@)?(?P<host>[\w.-]+)(:(?P<port>\d+))?/?(?P<dbname>[\w-]*)(?:\?(?P<parameters>.*))?$`)
	if !r.MatchString(dsn) {
		return nil, errors.New("wrong format. DSN must be in format: <scheme>://<user>:<password>@<host>:<port>/<dbname>?<parameters>")
	}
	matches := r.FindStringSubmatch(dsn)
	names := r.SubexpNames()

	m := make(map[string]string)
	for i, match := range matches {
		if i != 0 && names[i] != "" {
			m[names[i]] = match
		}
	}

	parameters := make(map[string]string)
	if m["parameters"] != "" {
		pairs := strings.Split(m["parameters"], "&")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) != 2 {
				return nil, errors.New("wrong format. Parameters must be in format: key=value")
			}
			parameters[kv[0]] = kv[1]
		}
	}

	d := dsntype{
		dsn:        dsn,
		scheme:     m["scheme"],
		user:       m["user"],
		password:   m["password"],
		host:       m["host"],
		port:       m["port"],
		dbname:     m["dbname"],
		parameters: parameters,
	}
	return &d, nil
}

func (d *dsntype) GetUser() string {
	return d.user
}

func (d *dsntype) GetPassword() string {
	return d.password
}

func (d *dsntype) GetHost() string {
	return d.host
}

func (d *dsntype) GetPort(defaultPort string) string {
	if d.port == "" {
		return defaultPort
	}
	return d.port
}

func (d *dsntype) GetPortInt(defaultPort int) int {
	portInt, err := strconv.Atoi(d.port)
	if err != nil {
		return defaultPort
	}
	return portInt
}

func (d *dsntype) GetParameter(parameter string) string {
	var value string
	if _, isKeyPresent := d.parameters[parameter]; isKeyPresent {
		value = d.parameters[parameter]
	}
	return value
}

func (d *dsntype) GetDBName() string {
	return d.dbname
}

func (d *dsntype) GetPostgresUri() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.GetHost(), d.GetPort("5432"), d.GetUser(), d.GetPassword(), d.GetDBName(), d.GetParameter("sslmode"))
	if d.GetParameter("search_path") != "" {
		psqlInfo = psqlInfo + fmt.Sprintf(" search_path=%s", d.GetParameter("search_path"))
	}
	return psqlInfo
}

func (d *dsntype) GetScheme() string {
	return d.scheme
}

func (d *dsntype) String() string {
	return d.dsn
}
