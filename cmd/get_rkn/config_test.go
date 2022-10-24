package main

import (
	"testing"
)

func Test_sanitizeConfLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		args   args
		wantS_ string
	}{
		{"", args{s: ""}, ""},
		{"", args{s: " "}, ""},
		{"", args{s: "    "}, ""},
		{"", args{s: "     #"}, ""},
		{"", args{s: "     #     "}, ""},
		{"", args{s: "#      "}, ""},
		{"", args{s: "#"}, ""},
		{"", args{s: "abc"}, "abc"},
		{"", args{s: "      abc"}, "abc"},
		{"", args{s: "abc      "}, "abc"},
		{"", args{s: "     abc    "}, "abc"},
		{"", args{s: "ab   cd"}, "ab   cd"},
		{"", args{s: "    ab     cd"}, "ab     cd"},
		{"", args{s: "    ab     cd       "}, "ab     cd"},
		{"", args{s: "#    ab     cd       "}, ""},
		{"", args{s: "    ab   #  cd       "}, "ab"},
		{"", args{s: "    ab     cd   #    "}, "ab     cd"},
		{"", args{s: "    ab     cd   #  asd  "}, "ab     cd"},
		{"", args{s: " #  asd  "}, ""},
		{"", args{s: "abc#asd"}, "abc"},
		{"", args{s: "    abc#asd"}, "abc"},
		{"", args{s: "    abc    #asd"}, "abc"},
		{"", args{s: "abc    #asd"}, "abc"},
		{"", args{s: "abc#     asd"}, "abc"},
		{"", args{s: "abc#  asd  #  adf"}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS_ := sanitizeConfLine(tt.args.s); gotS_ != tt.wantS_ {
				t.Errorf("sanitizeConfLine() = %v, want %v", gotS_, tt.wantS_)
			}
		})
	}
}
