// =====================================================================
//
// utils_test.go -
//
// Created by skywind on 2024/12/19
// Last Modified: 2024/12/19 16:24:25
//
// =====================================================================
package forward

import (
	"net"
	"testing"
)

func TestAddressSet(t *testing.T) {
	var dst, src net.TCPAddr
	src.IP = net.ParseIP("192.168.1.128")
	src.Port = 1234
	dst.IP = net.ParseIP("192.168.2.128")
	dst.Port = 5678
	AddressSet(&dst, &src)
	if dst.IP.String() != "192.168.1.128" {
		t.Fatalf("dst.IP = %s", dst.IP.String())
	}
	if dst.Port != 1234 {
		t.Fatalf("dst.Port = %d", dst.Port)
	}
}

func TestAddressResolve(t *testing.T) {
	addr := AddressResolve("192.168.1.118:1234")
	if addr == nil {
		t.Fatalf("resolve failed")
	}
	if addr.IP.String() != "192.168.1.118" {
		t.Fatalf("addr.IP = %s", addr.IP.String())
	}
	if addr.Port != 1234 {
		t.Fatalf("addr.Port = %d", addr.Port)
	}
}
