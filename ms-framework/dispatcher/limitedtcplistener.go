package dispatcher

import (
	"net"
	"sync"
)

type LimitedTcpListener struct {
	net.Listener
	sync.Mutex
	sem chan bool
	id  int
}

func InitLimitedTcpListener(count int, l net.Listener) LimitedTcpListener {
	sem := make(chan bool, count)
	for i := 0; i < count; i++ {
		sem <- true
	}

	return LimitedTcpListener{
		Listener: l,
		Mutex:    sync.Mutex{},
		sem:      sem,
		id:       123,
	}
}

func (l LimitedTcpListener) Addr() net.Addr { return l.Listener.Addr() }

func (l LimitedTcpListener) Close() error { return l.Listener.Close() }

func (l LimitedTcpListener) Accept() (net.Conn, error) {
	<-l.sem
	c, err := l.Listener.Accept()

	if err != nil {
		return nil, err
	}

	return &LimitedTcpConn{
		Conn: c,
		sem:  l.sem,
	}, nil
}

type LimitedTcpConn struct {
	net.Conn
	sem chan bool
}

func (c *LimitedTcpConn) Close() error {
	err := c.Conn.Close()
	c.sem <- true
	return err
}
