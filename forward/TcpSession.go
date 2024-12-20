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
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type TcpSession struct {
	name    string
	lock    sync.Mutex
	wg      sync.WaitGroup
	closing atomic.Bool
	local   *net.TCPConn
	remote  *net.TCPConn
	cin     *net.TCPConn
	cout    *net.TCPConn
	logger  *log.Logger
}

func NewTcpSession(name string, conn *net.TCPConn) *TcpSession {
	self := &TcpSession{
		name:   name,
		local:  conn,
		remote: nil,
		cin:    nil,
		cout:   nil,
		logger: nil,
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

func (self *TcpSession) SetLogger(logger *log.Logger) {
	self.logger = logger
}

func (self *TcpSession) SetRemote(addr *net.TCPAddr) error {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		if self.logger != nil {
			self.logger.Printf("[%s] remote connect error: %s", self.name, err)
		}
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
		if self.logger != nil {
			self.logger.Printf("[%s] input connect error: %s", self.name, err)
			self.logger.Printf("[%s] disable input", self.name)
		}
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
		if self.logger != nil {
			self.logger.Printf("[%s] output connect error: %s", self.name, err)
			self.logger.Printf("[%s] disable output", self.name)
		}
		return err
	}
	self.cout = conn
	return nil
}

func (self *TcpSession) blackhole(conn *net.TCPConn) {
	buf := make([]byte, 8192)
	for {
		if self.closing.Load() {
			break
		}
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		if n <= 0 {
			break
		}
	}
}

func (self *TcpSession) forward(src *net.TCPConn, dst *net.TCPConn, dup *net.TCPConn) {
	defer self.Close()
	if dup != nil {
		go self.blackhole(dup)
	}
	buf := make([]byte, 8192)
	lost := false
	for {
		n, err := src.Read(buf)
		if err != nil {
			if self.logger != nil {
				if src == self.local {
					self.logger.Printf("[%s] local read error: %s", self.name, err)
				} else {
					self.logger.Printf("[%s] remote read error: %s", self.name, err)
				}
			}
			break
		}
		if n <= 0 {
			if self.logger != nil {
				if src == self.local {
					self.logger.Printf("[%s] local read EOF", self.name)
				} else {
					self.logger.Printf("[%s] remote read EOF", self.name)
				}
			}
			break
		}
		_, err = dst.Write(buf[:n])
		if err != nil {
			if self.logger != nil {
				if src == self.local {
					self.logger.Printf("[%s] remote write error: %s", self.name, err)
				} else {
					self.logger.Printf("[%s] local write error: %s", self.name, err)
				}
			}
			break
		}
		if dup != nil && lost == false {
			_, err = dup.Write(buf[:n])
			if err != nil {
				if self.logger != nil {
					if src == self.local {
						self.logger.Printf("[%s] output write error: %s", self.name, err)
						self.logger.Printf("[%s] disable output", self.name)
					} else {
						self.logger.Printf("[%s] input write error: %s", self.name, err)
						self.logger.Printf("[%s] disable input", self.name)
					}
				}
				lost = true
			}
		}
	}
	self.wg.Done()
}

func (self *TcpSession) Run() {
	if self.logger != nil {
		self.logger.Printf("[%s] session started", self.name)
	}
	self.wg.Add(2)
	go self.forward(self.local, self.remote, self.cout)
	go self.forward(self.remote, self.local, self.cin)
	self.wg.Wait()
	if self.logger != nil {
		self.logger.Printf("[%s] session closed", self.name)
	}
}
