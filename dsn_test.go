package dsn

import (
	"reflect"
	"testing"
)

func TestDSN_GetUser(t *testing.T) {
	type fields struct {
		dsn string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "user postgres",
			fields: fields{
				dsn: "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			},
			want: "postgres",
		},
		{
			name: "wrong format",
			fields: fields{
				dsn: "postgres",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dsntype{
				dsn: tt.fields.dsn,
			}
			if got := d.GetUser(); got != tt.want {
				t.Errorf("DSN.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDSN_GetPostgresUri(t *testing.T) {
	type fields struct {
		dsn string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "wrong format",
			fields: fields{
				dsn: "toto",
			},
			want: "host= port=5432 user= password= dbname=toto sslmode=",
		},
		{
			name: "complete",
			fields: fields{
				dsn: "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			},
			want: "host=host port=5432 user=postgres password=password dbname=mydb sslmode=disable",
		},
		{
			name: "complete with search_path",
			fields: fields{
				dsn: "postgres://postgres:password@host:5432/mydb?sslmode=disable&search_path=namespace",
			},
			want: "host=host port=5432 user=postgres password=password dbname=mydb sslmode=disable search_path=namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dsntype{
				dsn: tt.fields.dsn,
			}
			if got := d.GetPostgresUri(); got != tt.want {
				t.Errorf("DSN.GetPostgresUri() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dsntype_GetPortInt(t *testing.T) {
	type fields struct {
		dsn string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "port string",
			fields: fields{
				dsn: "postgres://user:password@host:port/dbname",
			},
			want: 5432,
		},
		{
			name: "port int",
			fields: fields{
				dsn: "postgres://user:password@host:5432/dbname",
			},
			want: 5432,
		},
		{
			name: "no port defined",
			fields: fields{
				dsn: "postgres://user:password@host/dbname",
			},
			want: 5432,
		},
		{
			name: "no port defined",
			fields: fields{
				dsn: "host/dbname",
			},
			want: 5432,
		},
		{
			name: "pg://host:5433",
			fields: fields{
				dsn: "pg://host:5433",
			},
			want: 5433,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dsntype{
				dsn: tt.fields.dsn,
			}
			if got := d.GetPortInt(5432); got != tt.want {
				t.Errorf("dsntype.GetPortInt() = %v, want %v", got, tt.want)
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
		want    DSN
		wantErr bool
	}{
		{
			name: "ssh",
			args: args{
				dsn: "ssh://host",
			},
			want: &dsntype{
				dsn: "ssh://host",
			},
			wantErr: false,
		},
		{
			name: "pg",
			args: args{
				dsn: "pg://user:password@host:5432",
			},
			want: &dsntype{
				dsn: "pg://user:password@host:5432",
			},
			wantErr: false,
		},
		{
			name: "pg wrong port",
			args: args{
				dsn: "pg://user:password@host:wrongport",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("New() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_dsntype_GetPath(t *testing.T) {
	type fields struct {
		dsn string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty path",
			fields: fields{
				dsn: "redis://host",
			},
			want: "",
		},
		{
			name: "path=/0",
			fields: fields{
				dsn: "redis://host/0",
			},
			want: "/0",
		},
		{
			name: "/sdf/sdf/sdf",
			fields: fields{
				dsn: "/sdf/sdf/sdf",
			},
			want: "/sdf/sdf/sdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dsntype{
				dsn: tt.fields.dsn,
			}
			if got := d.GetPath(); got != tt.want {
				t.Errorf("dsntype.GetPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
