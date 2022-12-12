package tools

import (
	"net"
	"testing"
)

func Test_ipInSubnets(t *testing.T) {
	type args struct {
		ip      net.IP
		subnets []*net.IPNet
	}
	tests := []struct {
		name         string
		args         args
		wantContains bool
		wantSubnet   string
	}{
		{"empty", args{net.IP{}, []*net.IPNet{}}, false, ""},
		{"yes", args{net.IP{10, 10, 0, 15}, []*net.IPNet{
			{net.IP{10, 10, 1, 0}, net.IPMask{255, 255, 255, 0}},
			{net.IP{10, 10, 0, 0}, net.IPMask{255, 255, 255, 0}},
			{net.IP{192, 168, 22, 0}, net.IPMask{255, 255, 0, 0}},
		}}, true, "10.10.0.0/24"},
		{"no", args{net.IP{10, 10, 0, 15}, []*net.IPNet{
			{net.IP{10, 10, 1, 0}, net.IPMask{255, 255, 255, 0}},
			{net.IP{10, 10, 2, 0}, net.IPMask{255, 255, 255, 0}},
			{net.IP{192, 168, 22, 0}, net.IPMask{255, 255, 0, 0}},
		}}, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContains, gotSubnet := IpInSubnets(tt.args.ip, tt.args.subnets)
			if gotContains != tt.wantContains {
				t.Errorf("ipInSubnets() gotContains = %v, want %v", gotContains, tt.wantContains)
			}
			if gotSubnet != tt.wantSubnet {
				t.Errorf("ipInSubnets() gotSubnet = %v, want %v", gotSubnet, tt.wantSubnet)
			}
		})
	}
}
