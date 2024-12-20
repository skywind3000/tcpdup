// =====================================================================
//
// TcpForward.go -
//
// Created by skywind on 2024/12/19
// Last Modified: 2024/12/19 15:46:15
//
// =====================================================================
package forward

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// tcp forward
type TcpForward struct {
	lock     sync.Mutex
	listener *net.TCPListener
	closing  atomic.Bool
	wg       sync.WaitGroup
	srcAddr  *net.TCPAddr
	dstAddr  *net.TCPAddr
	input    string
	output   string
	logger   *log.Logger
}

func NewTcpForward() *TcpForward {
	self := &TcpForward{
		listener: nil,
		wg:       sync.WaitGroup{},
		srcAddr:  nil,
		dstAddr:  nil,
		logger:   nil,
		input:    "",
		output:   "",
	}
	self.closing.Store(false)
	return self
}

func (self *TcpForward) SetLogger(logger *log.Logger) {
	self.logger = logger
}

func (self *TcpForward) SetInput(address string) {
	self.input = address
}

func (self *TcpForward) SetOutput(address string) {
	self.output = address
}

func (self *TcpForward) Open(srcAddr *net.TCPAddr, dstAddr *net.TCPAddr) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
	self.srcAddr = AddressClone(srcAddr)
	self.dstAddr = AddressClone(dstAddr)
	self.closing.Store(false)
	listener, err := net.ListenTCP("tcp", srcAddr)
	if err != nil {
		return err
	}
	self.listener = listener
	self.wg.Add(1)
	notify := make(chan int)
	go func() {
		notify <- 1
		for self.closing.Load() == false {
			conn, err := self.listener.AcceptTCP()
			if err != nil {
				break
			}
			go self.handle(conn)
		}
		self.wg.Done()
	}()
	<-notify
	return nil
}

func (self *TcpForward) Wait() {
	self.wg.Wait()
}

func (self *TcpForward) shutdown() {
	self.closing.Store(true)
	if self.listener != nil {
		self.listener.Close()
	}
	self.listener = nil
	self.wg.Wait()
}

func (self *TcpForward) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
}

func (self *TcpForward) handle(conn *net.TCPConn) {
	defer conn.Close()
	name := conn.RemoteAddr().String()
	session := NewTcpSession(name, conn)
	session.SetLogger(self.logger)
	defer session.Close()
	if session.SetRemote(self.dstAddr) != nil {
		return
	}
	if session.SetInput(self.input) != nil {
		return
	}
	if session.SetOutput(self.output) != nil {
		return
	}
	session.Run()
}
