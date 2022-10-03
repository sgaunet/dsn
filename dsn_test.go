package dsn

import (
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
			d := &DSN{
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
			want: "host= port= user= password= dbname=toto sslmode=",
		},
		{
			name: "complete",
			fields: fields{
				dsn: "postgres://postgres:password@host:5432/mydb?sslmode=disable",
			},
			want: "host=host port=5432 user=postgres password=password dbname=mydb sslmode=disable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DSN{
				dsn: tt.fields.dsn,
			}
			if got := d.GetPostgresUri(); got != tt.want {
				t.Errorf("DSN.GetPostgresUri() = %v, want %v", got, tt.want)
			}
		})
	}
}
