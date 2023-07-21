package dsn_test

import (
	"testing"

	"github.com/sgaunet/dsn/pkg/dsn"
)

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
				if got := d.GetUser(); got != tt.want {
					t.Errorf("DSN.GetUser() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDSN_GetPostgresUri(t *testing.T) {
	tests := []struct {
		name      string
		dsnToTest string
		want      string
		wantErr   bool
	}{
		{
			name:      "wrong format",
			dsnToTest: "postgres",
			want:      "",
			wantErr:   true,
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
					if got := d.GetPostgresUri(); got != tt.want {
						t.Errorf("DSN.GetPostgresUri() = %v, want %v", got, tt.want)
					}
				}
			})
		})
	}
}

func Test_dsntype_GetPortInt(t *testing.T) {
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

func Test_dsntype_GetPort(t *testing.T) {
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
					t.Errorf("DSN.GetPortInt() = %v, want %v", got, tt.want)
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

func Test_dsntype_GetPath(t *testing.T) {
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

func Test_dsntype_GetScheme(t *testing.T) {
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
			wantErr:   true,
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

func Test_dsntype_GetHost(t *testing.T) {
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
			wantErr:   true,
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

func Test_dsntype_GetParameter(t *testing.T) {
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
			wantErr:   true,
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
			want:      "5",
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
					t.Errorf("DSN.GetHost() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
