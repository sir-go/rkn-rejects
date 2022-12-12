package main

import (
	"io/ioutil"
	"testing"
)

func TestTLSPayload_GetLenW(t *testing.T) {
	tests := []struct {
		name    string
		pl      TLSPayload
		want    int
		wantErr bool
	}{
		{"empty", TLSPayload{}, 0, true},
		{"tooSmall", TLSPayload{10, []byte{}, 16}, 0, true},
		{"ok", TLSPayload{6, []byte{10, 12, 2, 1, 0, 3}, 2}, 513, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pl.GetLenW()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLenW() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLenW() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTLSPayload_GetLenB(t *testing.T) {
	tests := []struct {
		name    string
		pl      TLSPayload
		want    int
		wantErr bool
	}{
		{"empty", TLSPayload{}, 0, true},
		{"tooSmall", TLSPayload{10, []byte{}, 16}, 0, true},
		{"ok", TLSPayload{6, []byte{10, 12, 2, 1, 0, 3}, 2}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pl.GetLenB()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLenB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLenB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTLSPayload_Skip(t *testing.T) {
	type args struct {
		pl *TLSPayload
		n  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantPos int
	}{
		{"empty", args{&TLSPayload{}, 0}, false, 0},
		{"tooSmall", args{&TLSPayload{}, 5}, true, 0},
		{"ok", args{&TLSPayload{20, nil, 3}, 5}, false, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.pl.Skip(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Skip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.pl.pos != tt.wantPos {
				t.Errorf("after Skip() pos = %v, want %v", tt.args.pl.pos, tt.wantPos)
			}
		})
	}
}

func TestTLSPayload_GetString(t *testing.T) {
	type args struct {
		pl *TLSPayload
		n  int
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr bool
	}{
		{"empty", args{&TLSPayload{}, 0}, "", false},
		{"tooSmall", args{&TLSPayload{}, 12}, "", true},
		{"ok", args{&TLSPayload{20, []byte("some string content"), 4}, 15},
			" string content", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.args.pl.GetString(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GetString() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetSNIForced(t *testing.T) {
	sniDumpVk, err := ioutil.ReadFile("../../testdata/sni_dump_vk.pcap")
	if err != nil {
		panic(err)
	}

	sniDumpTtnet, err := ioutil.ReadFile("../../testdata/sni_dump_ttnet.pcap")
	if err != nil {
		panic(err)
	}
	type args struct {
		d []byte
	}
	tests := []struct {
		name    string
		args    args
		wantSni string
		wantErr bool
	}{
		{"empty", args{nil}, "", true},
		{"vk", args{sniDumpVk}, "vk.com", false},
		{"ttnet", args{sniDumpTtnet}, "ttnet.ru", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSni, err := GetSNIForced(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSNIForced() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSni != tt.wantSni {
				t.Errorf("GetSNIForced() gotSni = %v, want %v", gotSni, tt.wantSni)
			}
		})
	}
}
