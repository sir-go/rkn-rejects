package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogBuff_add(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name        string
		args        args
		wantRecords []string
	}{
		{"empty", args{""}, []string{""}},
		{"msg", args{"some message"}, []string{"some message"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lB := &LogBuff{}
			lB.add(tt.args.msg)
			if !cmp.Equal(lB.records, tt.wantRecords) {
				t.Errorf("after LogBuff.add() records are %v, want %v", lB.records, tt.wantRecords)
			}
		})
	}
}

func TestLogBuff_flush(t *testing.T) {
	tests := []struct {
		name    string
		content []string
	}{
		{"empty", []string{""}},
		{"ok", []string{"some record", "another one record"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringsBuf := bytes.NewBufferString("")
			lB := &LogBuff{tt.content, stringsBuf}
			lB.flush()
			wantFlushed := strings.Join(tt.content, "\n") + "\n"
			if stringsBuf.String() != wantFlushed {
				t.Errorf("after LogBuff.flush() flushed data %s, want %s", stringsBuf.String(), wantFlushed)
			}
		})
	}
}
