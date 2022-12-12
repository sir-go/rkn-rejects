package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_dumpTo(t *testing.T) {
	type args struct {
		path string
		data interface{}
		decr string
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
	}{
		{"empty", args{"", nil, ""}, nil},
		{"ok", args{"/tmp/tmp-file-name", "some content", "description"},
			[]byte("some content")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dumpTo(&tt.args.path, tt.args.data, tt.args.decr)
			if tt.args.path != "" {
				content, err := ioutil.ReadFile(filepath.Clean(tt.args.path))
				if err != nil {
					t.Errorf("read tmp file error %v", err)
				}
				if !bytes.Equal(content, tt.wantContent) {
					t.Errorf("dumpTo() wrote an unexpected contant %v, want %v", content, tt.wantContent)
				}
			}
		})
	}
}
