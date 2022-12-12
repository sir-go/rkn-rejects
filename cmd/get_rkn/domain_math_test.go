package main

import (
	"testing"
)

func Test_isItUpDomainOf(t *testing.T) {
	type args struct {
		domain     string
		subDomains []string
	}
	tests := []struct {
		name          string
		args          args
		wantResult    bool
		wantSubdomain string
	}{
		{"empty", args{"", nil},
			false, ""},
		{"self", args{"example.com", []string{"example.com"}},
			true, "example.com"},
		{"child", args{"example.com", []string{"subdomain.example.com"}},
			true, "subdomain.example.com"},
		{"3rd", args{"subdomain.example.com", []string{"sub.subdomain.example.com"}},
			true, "sub.subdomain.example.com"},
		{"not", args{"subdomain.example.com", []string{"example.com"}},
			false, ""},
		{"many", args{"example.com", []string{
			"sub.another.com",
			"ech1.example.com",
			"bx.example.com",
		}},
			true, "ech1.example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotSubdomain := isItUpDomainOf(tt.args.domain, tt.args.subDomains)
			if gotResult != tt.wantResult || gotSubdomain != tt.wantSubdomain {
				t.Errorf("isItUpDomainOf() got = (%v, %v), want (%v, %v)",
					gotResult, gotSubdomain, tt.wantResult, tt.wantSubdomain)
			}
		})
	}
}

func Test_sanitizeDomain(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		{"empty", args{""}, ""},
		{"www", args{"www.ya.ru"}, "ya.ru"},
		{"wildcard", args{"*.maps.ya.ru"}, "maps.ya.ru"},
		{"params", args{"www.maps.ya.ru/page?params=4&l=3"}, "maps.ya.ru"},
		{"anchors", args{"*.maps.ya.ru/#:page"}, "maps.ya.ru"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := sanitizeDomain(tt.args.d); gotS != tt.wantS {
				t.Errorf("sanitizeDomain() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func Test_forcePDecode(t *testing.T) {
	type args struct {
		d string
	}
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		{"empty", args{""}, ""},
		{"national", args{"xn--e1afmkfd.xn--p1ai"}, "пример.рф"},
		{"eng", args{"example.com"}, "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := forcePDecode(tt.args.d); gotS != tt.wantS {
				t.Errorf("forcePDecode() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func Test_forcePUnicode(t *testing.T) {
	type args struct {
		d string
	}
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		{"empty", args{""}, ""},
		{"national", args{"пример.рф"}, "xn--e1afmkfd.xn--p1ai"},
		{"eng", args{"example.com"}, "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := forcePUnicode(tt.args.d); gotS != tt.wantS {
				t.Errorf("forcePUnicode() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func Test_isASCII(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{""}, true},
		{"yes", args{"latin string"}, true},
		{"no", args{"latin строка"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isASCII(tt.args.s); got != tt.want {
				t.Errorf("isASCII() = %v, want %v", got, tt.want)
			}
		})
	}
}
