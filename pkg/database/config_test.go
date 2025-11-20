package database

import (
	"testing"
)

func TestConfig_ConnectionURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			name:   "nil",
			config: nil,
			want:   "",
		},
		{
			name: "host",
			config: &Config{
				Host: "myhost",
			},
			want: "postgres://myhost:5432",
		},
		{
			name: "host_port",
			config: &Config{
				Host: "myhost",
				Port: "1234",
			},
			want: "postgres://myhost:1234",
		},
		{
			name: "basic_auth",
			config: &Config{
				User:     "myuser",
				Password: "mypass",
			},
			want: "postgres://myuser:mypass@localhost:5432",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg := tc.config
			if got, want := cfg.ConnectionURL(), tc.want; got != want {
				t.Errorf("expected %q to be %q", got, want)
			}
		})
	}

}
