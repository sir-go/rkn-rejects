package main

import (
	"net"
	"testing"
)

func Test_ipInSubnets(t *testing.T) {
	type args struct {
		ip      net.IP
		subnets []Net
	}
	tests := []struct {
		name       string
		args       args
		wantFound  bool
		wantSubnet string
	}{
		{
			name: "empty",
			args: args{
				ip:      nil,
				subnets: nil,
			},
			wantFound:  false,
			wantSubnet: "",
		},
		{
			name: "found",
			args: args{
				ip: net.IP{10, 10, 0, 15},
				subnets: []Net{
					{IPNet: &net.IPNet{IP: net.IP{10, 10, 10, 0}, Mask: net.IPMask{255, 255, 255, 0}}},
					{IPNet: &net.IPNet{IP: net.IP{10, 10, 0, 0}, Mask: net.IPMask{255, 255, 255, 0}}},
					{IPNet: &net.IPNet{IP: net.IP{192, 168, 201, 0}, Mask: net.IPMask{255, 255, 255, 252}}},
				},
			},
			wantFound:  true,
			wantSubnet: "10.10.0.0/24",
		},
		{
			name: "not-found",
			args: args{
				ip: net.IP{10, 10, 0, 15},
				subnets: []Net{
					{IPNet: &net.IPNet{IP: net.IP{10, 10, 10, 0}, Mask: net.IPMask{255, 255, 255, 0}}},
					{IPNet: &net.IPNet{IP: net.IP{172, 16, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}}},
					{IPNet: &net.IPNet{IP: net.IP{192, 168, 201, 0}, Mask: net.IPMask{255, 255, 255, 252}}},
				},
			},
			wantFound:  false,
			wantSubnet: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFound, gotSubnet := ipInSubnets(tt.args.ip, tt.args.subnets)
			if gotFound != tt.wantFound || gotSubnet != tt.wantSubnet {
				t.Errorf("ipInSubnets() got = (%v, %v), want (%v, %v)",
					gotFound, gotSubnet, tt.wantFound, tt.wantSubnet)
			}
		})
	}
}
