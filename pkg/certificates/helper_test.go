package certificates

import (
	"testing"
)

func Test_Sanitize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"single cert",
			args{s: "-----BEGIN CERTIFICATE-----\n1\n-----END CERTIFICATE-----\n\n"},
			"-----BEGIN CERTIFICATE-----\n1\n-----END CERTIFICATE-----\n",
		},
		{
			"multiple certs without extra newline",
			args{s: "-----BEGIN CERTIFICATE-----\n1\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\n2\n-----END CERTIFICATE-----\n\n"},
			"-----BEGIN CERTIFICATE-----\n1\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\n2\n-----END CERTIFICATE-----\n",
		},
		{
			"multiple certs with extra newline",
			args{s: "-----BEGIN CERTIFICATE-----\n1\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\n2\n-----END CERTIFICATE-----\n\n"},
			"-----BEGIN CERTIFICATE-----\n1\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\n2\n-----END CERTIFICATE-----\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sanitize(tt.args.s); got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
