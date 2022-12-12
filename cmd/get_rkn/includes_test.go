package main

import (
	"testing"
)

func Test_sanitizeConfLine(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{"", ""},
		{" ", ""},
		{"    ", ""},
		{"     #", ""},
		{"     #     ", ""},
		{"#      ", ""},
		{"#", ""},
		{"abc", "abc"},
		{"      abc", "abc"},
		{"abc      ", "abc"},
		{"     abc    ", "abc"},
		{"ab   cd", "ab   cd"},
		{"    ab     cd", "ab     cd"},
		{"    ab     cd       ", "ab     cd"},
		{"#    ab     cd       ", ""},
		{"    ab   #  cd       ", "ab"},
		{"    ab     cd   #    ", "ab     cd"},
		{"    ab     cd   #  asd  ", "ab     cd"},
		{" #  asd  ", ""},
		{"abc#asd", "abc"},
		{"    abc#asd", "abc"},
		{"    abc    #asd", "abc"},
		{"abc    #asd", "abc"},
		{"abc#     asd", "abc"},
		{"abc#  asd  #  adf", "abc"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := sanitizeConfLine(tt.arg); got != tt.want {
				t.Errorf("sanitizeConfLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
