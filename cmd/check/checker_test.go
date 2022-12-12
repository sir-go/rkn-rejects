package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// checks if dump file is created and dumped data equals for verdict's raw data
func Test_verdict_dump(t *testing.T) {
	type args struct {
		vDir string
	}
	tests := []struct {
		name    string
		args    args
		verdict verdict
	}{
		{"e2e",
			args{"_verdicts_dump_test"},
			verdict{
				true,
				"some-verdict-hash",
				"some-target",
				[]byte("some response content")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				tmpDirName string
				err        error
			)
			if tmpDirName, err = os.MkdirTemp("", tt.args.vDir); err != nil {
				t.Errorf("can't make a directory %s. %v", tt.args.vDir, err)
			}
			tt.verdict.dump(tmpDirName)
			dumpFile := filepath.Clean(filepath.Join(tmpDirName, tt.verdict.hash))
			verdictBytes, err := ioutil.ReadFile(dumpFile)
			if err != nil {
				t.Errorf("can't read the dump file %s, %v", dumpFile, err)
			}
			if !bytes.Equal(verdictBytes, tt.verdict.raw) {
				t.Errorf("dumped data %s, want %s", verdictBytes, tt.verdict.raw)
			}
		})
	}
}

func Test_check(t *testing.T) {
	type args struct {
		target  string
		timeout time.Duration
	}
	tests := []struct {
		name  string
		args  args
		wantV verdict
	}{
		{"google.com", args{"google.com", 10 * time.Second}, verdict{
			opened: true,
			target: "google.com",
			raw:    []byte("200 OK\n"),
		}},
		{"example.com", args{"example.com", 10 * time.Second}, verdict{
			opened: true,
			target: "example.com",
			raw:    []byte("200 OK\n"),
		}},
		{"non-exist", args{"some-non-exist-target-url.es", 10 * time.Second}, verdict{
			opened: false,
			target: "",
			raw:    nil,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV := check(tt.args.target, tt.args.timeout)
			if gotV.target != tt.wantV.target {
				t.Errorf("check(); verdict.target = %v, want %v", gotV.target, tt.wantV.target)
			}
			if gotV.opened != tt.wantV.opened {
				t.Errorf("check(); verdict.opened = %v, want %v", gotV.opened, tt.wantV.opened)
			}
			if !bytes.HasPrefix(gotV.raw, tt.wantV.raw) {
				t.Errorf("check(); verdict.raw begins with %s, want %s", gotV.raw, tt.wantV.raw)
			}
		})
	}
}
