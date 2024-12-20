// =====================================================================
//
// TcpSession.go -
//
// Created by skywind on 2024/12/20
// Last Modified: 2024/12/20 10:45:10
//
// =====================================================================
package forward

import (
	"net"
	"sync"
	"sync/atomic"
)

type TcpSession struct {
	lock    sync.Mutex
	closing atomic.Bool
	local   *net.TCPConn
	remote  *net.TCPConn
	cin     *net.TCPConn
	cout    *net.TCPConn
}

func NewTcpSession(conn *net.TCPConn) *TcpSession {
	self := &TcpSession{
		local:  conn,
		remote: nil,
		cin:    nil,
		cout:   nil,
	}
	self.closing.Store(false)
	return self
}

func (self *TcpSession) shutdown() {
	self.closing.Store(true)
	if self.local != nil {
		self.local.Close()
		self.local = nil
	}
	if self.remote != nil {
		self.remote.Close()
		self.remote = nil
	}
	if self.cin != nil {
		self.cin.Close()
		self.cin = nil
	}
	if self.cout != nil {
		self.cout.Close()
		self.cout = nil
	}
}

func (self *TcpSession) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
}

func (self *TcpSession) IsClosing() bool {
	return self.closing.Load()
}

func (self *TcpSession) SetRemote(addr *net.TCPAddr) error {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	self.remote = conn
	return nil
}

func (self *TcpSession) SetInput(addr string) error {
	self.cin = nil
	if addr == "" {
		return nil
	}
	a := AddressResolve(addr)
	if a == nil {
		return nil
	}
	conn, err := net.DialTCP("tcp", nil, a)
	if err != nil {
		return err
	}
	self.cin = conn
	return nil
}

func (self *TcpSession) SetOutput(addr string) error {
	self.cout = nil
	if addr == "" {
		return nil
	}
	a := AddressResolve(addr)
	if a == nil {
		return nil
	}
	conn, err := net.DialTCP("tcp", nil, a)
	if err != nil {
		return err
	}
	self.cout = conn
	return nil
}

func (self *TcpSession) Start() {
}
