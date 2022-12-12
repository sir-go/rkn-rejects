package main

import (
	"regexp"
	"testing"
)

func Test_snatchRegexp(t *testing.T) {
	type args struct {
		b   []byte
		pos int
	}
	tests := []struct {
		name string
		args args
		want map[*regexp.Regexp]string
	}{
		{"empty", args{[]byte(""), 0}, map[*regexp.Regexp]string{
			ReHTTP:            "",
			ReHostname:        "",
			ReLooksLikeDomain: "",
		}},
		{"domain", args{[]byte("host: example.com "), 0}, map[*regexp.Regexp]string{
			ReHTTP:            "",
			ReHostname:        "host: example.com ",
			ReLooksLikeDomain: "example.com",
		}},
		{"host", args{[]byte("host: example.com "), 1}, map[*regexp.Regexp]string{
			ReHTTP:            "",
			ReHostname:        "example.com",
			ReLooksLikeDomain: "example",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for rxp, wantString := range tt.want {
				if got := snatchRegexp(tt.args.b, rxp, tt.args.pos); got != wantString {
					t.Errorf("snatchRegexp(%s, %v, %v) = %v, want %v",
						tt.args.b, rxp, tt.args.pos, got, wantString)
				}
			}
		})
	}
}
