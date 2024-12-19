// =====================================================================
//
// utils.go -
//
// Created by skywind on 2024/12/19
// Last Modified: 2024/12/19 16:22:18
//
// =====================================================================
package forward

import (
	"net"
	"strings"
)

func AddressSet(dst *net.TCPAddr, src *net.TCPAddr) *net.TCPAddr {
	if len(dst.IP) != len(src.IP) {
		dst.IP = make(net.IP, len(src.IP))
	}
	copy(dst.IP, src.IP)
	dst.Port = src.Port
	dst.Zone = src.Zone
	return dst
}

func AddressClone(addr *net.TCPAddr) *net.TCPAddr {
	naddr := &net.TCPAddr{
		IP:   make(net.IP, len(addr.IP)),
		Port: addr.Port,
		Zone: addr.Zone,
	}
	copy(naddr.IP, addr.IP)
	return naddr
}

func AddressResolve(address string) *net.TCPAddr {
	if !strings.Contains(address, ":") {
		address = ":" + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil
	}
	return addr
}

func AddressString(addr *net.TCPAddr) string {
	return addr.String()
}
