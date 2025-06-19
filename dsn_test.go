package dsn_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/ghcr.io/sgaunet/dsn/v3"
)

func FuzzGetHost(f *testing.F) {
	testcases := []string{"postgres://postgres:password@postgres-server:5432/postgres?sslmode=disable",
		"postgres://postgres:password@host:5432/mydb?sslmode=disable",
		"hostname",
	}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		d, err := dsn.New(orig)
		if err != nil {
			t.Skip()
		}
		host := d.GetHost()
		scheme := d.GetScheme()
		if !utf8.ValidString(host) {
			t.Errorf("GetHost produced invalid UTF-8 string %q", host)
		}
		if !utf8.ValidString(scheme) {
			t.Errorf("GetScheme produced invalid UTF-8 string %q", scheme)
		}
	})
}

func TestDSN_GetUser(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "user postgres",
			dsnToTest: "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			want:      "postgres",
			wantErr:   false,
		},
		{
			name:      "wrong format",
			dsnToTest: "postgres",
			want:      "",
			wantErr:   false,
		},
		{
			name:      "password with special characters",
			dsnToTest: "postgres://user-name:pas&!sword-with-hyphens@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			want:      "user-name",
			wantErr:   false,
		},
		{
			name:      "field with -",
			dsnToTest: "postgres://user-name:password-with-hyphens@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			want:      "user-name",
			wantErr:   false,
		},
		{
			name:      "...",
			dsnToTest: "postgres://postgres:password@postgres-server:5432/postgres?sslmode=disable",
			want:      "postgres",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetUser(); got != tt.want {
					t.Errorf("DSN.GetUser() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDSN_GetPostgresURI(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "wrong format",
			dsnToTest: "postgres",
			want:      "host=postgres port=5432 user= password= dbname= sslmode=",
			wantErr:   false,
		},
		{
			name:      "complete",
			dsnToTest: "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			want:      "host=host port=5432 user=postgres password=password dbname=mydb sslmode=disable",
			wantErr:   false,
		},
		{
			name:      "complete with search_path",
			dsnToTest: "postgres://postgres:password@host:5432/mydb?sslmode=disable&search_path=namespace",
			want:      "host=host port=5432 user=postgres password=password dbname=mydb sslmode=disable search_path=namespace",
			wantErr:   false,
		},
		{
			name:      "...",
			dsnToTest: "postgres://postgres:password@postgres-server:5432/postgres?sslmode=disable",
			want:      "host=postgres-server port=5432 user=postgres password=password dbname=postgres sslmode=disable",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				d, err := dsn.New(tt.dsnToTest)
				if (err != nil) != tt.wantErr {
					t.Errorf("New(), wantErr %v", tt.wantErr)
					return
				}
				if tt.wantErr && d != nil {
					t.Errorf("expected d to be nil")
				}
				if !tt.wantErr {
					if got := d.GetPostgresURI(); got != tt.want {
						t.Errorf("DSN.GetPostgresURI() = %v, want %v", got, tt.want)
					}
				}
			})
		})
	}
}

func Test_DSN_GetPortInt(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      int
		wantErr   bool
	}{
		{
			name:      "port string",
			dsnToTest: "postgres://user:password@host:port/dbname",
			want:      5432,
			wantErr:   true,
		},
		{
			name:      "port int",
			dsnToTest: "postgres://user:password@host:5432/dbname",
			want:      5432,
			wantErr:   false,
		},
		{
			name:      "no port defined",
			dsnToTest: "postgres://user:password@host/dbname",
			want:      5432,
			wantErr:   false,
		},
		{
			name:      "pg://host:5433",
			dsnToTest: "pg://host:5433",
			want:      5433,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetPortInt(5432); got != tt.want {
					t.Errorf("DSN.GetPortInt() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_DSN_GetPort(t *testing.T) {
	defaultPort := "5432"
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "port int",
			dsnToTest: "postgres://user:password@host:5432/dbname",
			want:      "5432",
			wantErr:   false,
		},
		{
			name:      "no port defined",
			dsnToTest: "postgres://user:password@host/dbname",
			want:      "5432",
			wantErr:   false,
		},
		{
			name:      "pg://host:5433",
			dsnToTest: "pg://host:5433",
			want:      "5433",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetPort(defaultPort); got != tt.want {
					t.Errorf("DSN.GetPort() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ssh",
			args: args{
				dsn: "ssh://host",
			},
			wantErr: false,
		},
		{
			name: "pg",
			args: args{
				dsn: "pg://user:password@host:5432",
			},
			wantErr: false,
		},
		{
			name: "pg wrong port",
			args: args{
				dsn: "pg://user:password@host:port",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
		})
	}
}

func Test_DSN_GetPath(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "empty path",
			dsnToTest: "redis://host",
			want:      "",
			wantErr:   false,
		},
		{
			name:      "path=/0",
			dsnToTest: "redis://host/0",
			want:      "0",
			wantErr:   false,
		},
		{
			name:      "/sdf/sdf/sdf",
			dsnToTest: "/sdf/sdf/sdf",
			want:      "/sdf/sdf/sdf",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetDBName(); got != tt.want {
					t.Errorf("DSN.GetPath() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_DSN_GetScheme(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "empty scheme",
			dsnToTest: "host",
			want:      "",
			wantErr:   false,
		},
		{
			name:      "http scheme",
			dsnToTest: "http://url.com",
			want:      "http",
			wantErr:   false,
		},
		{
			name:      "https scheme",
			dsnToTest: "https://url.com",
			want:      "https",
			wantErr:   false,
		},
		{
			name:      "pg scheme",
			dsnToTest: "pg://user:password@sdfsdf",
			want:      "pg",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetScheme(); got != tt.want {
					t.Errorf("DSN.GetScheme() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_DSN_GetHost(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "empty scheme",
			dsnToTest: "host",
			want:      "host",
			wantErr:   false,
		},
		{
			name:      "http scheme",
			dsnToTest: "http://url.com",
			want:      "url.com",
			wantErr:   false,
		},
		{
			name:      "https scheme",
			dsnToTest: "https://url.com",
			want:      "url.com",
			wantErr:   false,
		},
		{
			name:      "pg scheme",
			dsnToTest: "pg://user:password@sdfsdf",
			want:      "sdfsdf",
			wantErr:   false,
		},
		{
			name:      "...",
			dsnToTest: "postgres://postgres:password@postgres-server:5432/postgres?sslmode=disable",
			want:      "postgres-server",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetHost(); got != tt.want {
					t.Errorf("DSN.GetHost() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_DSN_GetParameter(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		parameter string
		want      string
		wantErr   bool
	}{
		{
			name:      "empty parameter",
			dsnToTest: "host",
			parameter: "test",
			want:      "",
			wantErr:   false,
		},
		{
			name:      "empty parameter with good dsn",
			dsnToTest: "http://url.com",
			parameter: "test",
			want:      "",
			wantErr:   false,
		},
		{
			name:      "simple parameter",
			dsnToTest: "https://url.com/db?test=1",
			parameter: "test",
			want:      "1",
			wantErr:   false,
		},
		{
			name:      "double parameter",
			dsnToTest: "pg://user:password@sdfsdf?test=5&test=sdfsdf",
			parameter: "test",
			want:      "sdfsdf",
			wantErr:   false,
		},
		{
			name:      "multiple parameter",
			dsnToTest: "pg://user:password@sdfsdf?test5=5&test3=sdfsdf",
			parameter: "test5",
			want:      "5",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetParameter(tt.parameter); got != tt.want {
					t.Errorf("DSN.GetParameter() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "user postgres",
			dsnToTest: "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			want:      "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			wantErr:   false,
		},
		{
			name:      "wrong format",
			dsnToTest: "postgres",
			want:      "postgres",
			wantErr:   false,
		},
		{
			name:      "password with special characters",
			dsnToTest: "postgres://user-name:pas&!sword-with-hyphens@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			want:      "postgres://user-name:pas&!sword-with-hyphens@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			wantErr:   false,
		},
		{
			name:      "field with -",
			dsnToTest: "postgres://user-name:password-with-hyphens@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			want:      "postgres://user-name:password-with-hyphens@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.String(); got != tt.want {
					t.Errorf("DSN.String() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_DSN_GetPassword(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "simple password",
			dsnToTest: "postgres://user:password@host:5432/dbname",
			want:      "password",
			wantErr:   false,
		},
		{
			name:      "no port defined",
			dsnToTest: "postgres://user:password@host/dbname",
			want:      "password",
			wantErr:   false,
		},
		{
			name:      "no password",
			dsnToTest: "pg://host:5433",
			want:      "",
			wantErr:   false,
		},
		{
			name:      "complex password",
			dsnToTest: "pg://user:p%40ssw0rd!@host:5433",
			want:      "p@ssw0rd!",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v, got %v", tt.wantErr, err)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				if got := d.GetPassword(); got != tt.want {
					t.Errorf("DSN.GetPassword() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetParameters(t *testing.T) {
	// Create a DSN with multiple parameters
	dsnString := "postgres://user:password@host:5432/mydb?sslmode=disable&connect_timeout=10&search_path=public"
	d, err := dsn.New(dsnString)
	if err != nil {
		t.Fatalf("Failed to create DSN: %v", err)
	}

	// Get the parameters map
	params := d.GetParameters()

	// Verify expected parameters exist with correct values
	expected := map[string]string{
		"sslmode":         "disable",
		"connect_timeout": "10",
		"search_path":     "public",
	}

	for k, v := range expected {
		if params[k] != v {
			t.Errorf("GetParameters()[%q] = %q, want %q", k, params[k], v)
		}
	}

	// Verify map returned is a copy by modifying it
	params["new_param"] = "test"

	// Get parameters again and verify the original wasn't modified
	newParams := d.GetParameters()
	if _, exists := newParams["new_param"]; exists {
		t.Error("GetParameters() should return a copy, not a reference to the internal map")
	}
}

func TestSetParameter(t *testing.T) {
	tests := []struct {
		name       string
		dsnToTest  string
		paramKey   string
		paramValue string
		want       string
		wantErr    bool
	}{
		{
			name:       "add new parameter",
			dsnToTest:  "postgres://host:5432/mydb",
			paramKey:   "sslmode",
			paramValue: "disable",
			want:       "disable",
			wantErr:    false,
		},
		{
			name:       "update existing parameter",
			dsnToTest:  "postgres://host:5432/mydb?sslmode=require",
			paramKey:   "sslmode",
			paramValue: "disable",
			want:       "disable",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v, got %v", tt.wantErr, err)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				// Set the parameter
				d = d.SetParameter(tt.paramKey, tt.paramValue)

				// Verify it was set correctly
				if got := d.GetParameter(tt.paramKey); got != tt.want {
					t.Errorf("After SetParameter, DSN.GetParameter(%q) = %v, want %v", tt.paramKey, got, tt.want)
				}

				// Check if parameters map contains expected value
				params := d.GetParameters()
				if params[tt.paramKey] != tt.want {
					t.Errorf("Parameters map[%q] = %v, want %v", tt.paramKey, params[tt.paramKey], tt.want)
				}
			}
		})
	}
}

func TestStringRedacted(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "no password",
			dsnToTest: "postgres://host:5432/mydb",
			want:      "postgres://host:5432/mydb",
			wantErr:   false,
		},
		{
			name:      "with password",
			dsnToTest: "postgres://user:password@host:5432/mydb?sslmode=disable",
			want:      "postgres://user:****@host:5432/mydb?sslmode=disable",
			wantErr:   false,
		},
		{
			name:      "with complex password",
			dsnToTest: "postgres://user-name:p@ssw0rd!@host-suffix.domain.com:5432/mydb-3?sslmode=require&connect_timeout=10s",
			want:      "postgres://user-name:****@host-suffix.domain.com:5432/mydb-3?", // We'll check parameters separately
			wantErr:   false,
		},
		{
			name:      "with no scheme",
			dsnToTest: "user:password@host",
			want:      "user:****@host",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := dsn.New(tt.dsnToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(), wantErr %v, got %v", tt.wantErr, err)
				return
			}
			if tt.wantErr && d != nil {
				t.Errorf("expected d to be nil")
			}
			if !tt.wantErr {
				got := d.StringRedacted()
				if tt.name == "with complex password" {
					// For complex password case, check the base URL and parameters separately
					if !strings.HasPrefix(got, tt.want) {
						t.Errorf("DSN.StringRedacted() = %v, should start with %v", got, tt.want)
					}
					// Also check that both parameters exist
					if !strings.Contains(got, "sslmode=require") {
						t.Errorf("DSN.StringRedacted() missing parameter sslmode=require")
					}
					if !strings.Contains(got, "connect_timeout=10s") {
						t.Errorf("DSN.StringRedacted() missing parameter connect_timeout=10s")
					}
				} else if got != tt.want {
					t.Errorf("DSN.StringRedacted() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
