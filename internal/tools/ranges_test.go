package tools

import (
	"reflect"
	"testing"
)

func Test_parseRange(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantSlice []uint16
		wantErr   bool
	}{
		{"single", args{"34"}, []uint16{34}, false},
		{"singleErr", args{"3f4"}, []uint16{}, true},
		{"range", args{"34-36"}, []uint16{34, 35, 36}, false},
		{"rangeErr", args{"34-3-6"}, []uint16{}, false},
	}
	for _, tt := range tests {
		a := []uint16{}
		t.Run(tt.name, func(t *testing.T) {
			err := ParseRange(tt.args.s, &a)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(a, tt.wantSlice) {
				t.Errorf("after parseRange() given slice contains %v, want %v ", a, tt.wantSlice)
			}
		})
	}
}
