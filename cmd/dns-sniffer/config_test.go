package main

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
	}{
		{"single", args{"34"}, []uint16{34}},
		{"range", args{"34-36"}, []uint16{34, 35, 36}},
	}
	for _, tt := range tests {
		a := []uint16{}
		t.Run(tt.name, func(t *testing.T) {
			parseRange(tt.args.s, &a)
			if !reflect.DeepEqual(a, tt.wantSlice) {
				t.Errorf("after parseRange() given slice contains %v, want %v ", a, tt.wantSlice)
			}
		})
	}
}
