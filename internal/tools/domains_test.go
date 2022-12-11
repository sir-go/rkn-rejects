package tools

import (
	"reflect"
	"testing"
)

func Test_getUpperDomains(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{"empty", args{""}, []string{""}},
		{"1lvl", args{"com"}, []string{"com"}},
		{"2lvl", args{"uk.com"}, []string{"uk.com"}},
		{"3lvl", args{"gov.uk.com"}, []string{"uk.com", "gov.uk.com"}},
		{"4lvl", args{"main.gov.uk.com"}, []string{"uk.com", "gov.uk.com", "main.gov.uk.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := GetUpperDomains(tt.args.d); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getUpperDomains() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
