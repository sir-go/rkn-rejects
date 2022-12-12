package main

import (
	"net"
	"reflect"
	"testing"
)

func TestNet_parse(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *net.IPNet
	}{
		{"empty", args{""}, true, nil},
		{"ok-Zero", args{"0.0.0.0/0"}, false,
			&net.IPNet{IP: net.IP{0, 0, 0, 0}, Mask: net.IPMask{0, 0, 0, 0}}},
		{"ok-Net", args{"192.168.5.0/30"}, false,
			&net.IPNet{IP: net.IP{192, 168, 5, 0}, Mask: net.IPMask{255, 255, 255, 252}}},
		{"ok-Host", args{"10.10.6.36/32"}, false,
			&net.IPNet{IP: net.IP{10, 10, 6, 36}, Mask: net.IPMask{255, 255, 255, 255}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipN := &Net{}
			if err := ipN.parse(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(ipN.IPNet.String(), tt.want.String()) {
				t.Errorf("after parse() ipNet = %v, want %v", ipN.IPNet, tt.want)
			}
		})
	}
}
