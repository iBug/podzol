package server

import (
	"bufio"
	"net"
	"sync"
)

type HTTPServer struct {
	s *Server
}

const BUFSIZE = 8192

var pool = &sync.Pool{
	New: func() any {
		return make([]byte, 0, BUFSIZE)
	},
}

func getBuffer() []byte {
	return pool.Get().([]byte)
}

func putBuffer(buf []byte) {
	pool.Put(buf[:0])
}

func closeConn(conn *net.TCPConn) {
	conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nServer: ustclug/podzol\r\n\r\nMissing Host header\n"))
	conn.Close()
}

// Create an HTTPServer from a Server.
func (s *Server) HTTPServer() *HTTPServer {
	return &HTTPServer{s}
}

func (s *HTTPServer) Handle(conn *net.TCPConn) {
	// Note: This is still very simple and contains lots of allocations.
	r := bufio.NewReaderSize(conn, BUFSIZE)
	buf := getBuffer()
	defer putBuffer(buf)
	defer r.Discard(0)
}
