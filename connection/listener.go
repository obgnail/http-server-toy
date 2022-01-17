package connection

import (
	"fmt"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

type Listener struct {
	addr        net.Addr
	tcpListener *net.TCPListener
	connChan    chan *Conn
	once        sync.Once
}

func (l *Listener) Close() {
	l.once.Do(l.close)
}

func (l *Listener) close() {
	l.tcpListener.Close()
}

func (l *Listener) startListen() {
	if l.tcpListener == nil {
		panic("tcpListener == nil")
	}
	log.Info("start listen:", l.addr)
	for {
		conn, err := l.tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		log.Infof("get remote conn: %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
		c := NewConn(conn)
		l.connChan <- c
	}
}

func NewListener(bindAddr string, bindPort int64) (listener *Listener, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return listener, errors.Trace(err)
	}
	listener = &Listener{
		addr:        tcpListener.Addr(),
		tcpListener: tcpListener,
		connChan:    make(chan *Conn, 1024),
	}
	go listener.startListen()
	return listener, nil
}

func (l *Listener) GetConn() (conn *Conn, err error) {
	var ok bool
	conn, ok = <-l.connChan
	if !ok {
		return conn, fmt.Errorf("channel close")
	}
	return conn, nil
}
