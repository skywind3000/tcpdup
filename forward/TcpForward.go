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
	"io"
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
	dst, err := net.DialTCP("tcp", nil, self.dstAddr)
	if err != nil {
		return
	}
	defer dst.Close()
	copyDie := make(chan int)
	go func() {
		io.Copy(conn, dst)
		copyDie <- 1
	}()
	io.Copy(dst, conn)
	<-copyDie
}
