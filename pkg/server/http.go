package server

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"strings"
	"sync"
)

type HTTPServer struct {
	s *Server
}

const BUFSIZE = 8192

var pool = &sync.Pool{
	New: func() any {
		buf := make([]byte, 0, BUFSIZE)
		return &buf
	},
}

func getBuffer() []byte {
	return *pool.Get().(*[]byte)
}

func putBuffer(buf []byte) {
	buf = buf[:0]
	pool.Put(&buf)
}

func closeConn(conn *net.TCPConn) {
	conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nServer: ustclug/podzol\r\n\r\nMissing Host header\n"))
}

// Create an HTTPServer from a Server.
func (s *Server) HTTPServer() *HTTPServer {
	return &HTTPServer{s}
}

func (s *HTTPServer) Handle(conn *net.TCPConn) {
	defer conn.Close()

	// Note: This is still very simple and contains lots of allocations.
	r := bufio.NewReaderSize(conn, BUFSIZE)
	buf := getBuffer()
	bufi := 0
	defer putBuffer(buf)
	upstreamAddr := ""
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			break
		}
		if bufi+len(line) > len(buf) {
			// too large
			closeConn(conn)
			return
		}

		if len(line) > 5 && bytes.EqualFold(line[:5], []byte("Host:")) {
			hostname := string(bytes.TrimSpace(line[5:]))
			hostname = strings.SplitN(hostname, ".", 2)[0]
			upstreamAddr, err = s.s.docker.GetIP(context.TODO(), hostname)
			if err != nil {
				closeConn(conn)
				return
			}
			break
		}
		if len(line) < 2 || bytes.Equal(line, []byte("\r\n")) {
			// all headers read
			break
		}

		// accumulate to buf
		bufi += copy(buf[bufi:], line)
	}
	if upstreamAddr == "" {
		closeConn(conn)
		return
	}

	// Connect to upstream
	upstreamConnTemp, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		closeConn(conn)
		return
	}
	upstreamConn := upstreamConnTemp.(*net.TCPConn)
	defer upstreamConn.Close()
	defer conn.Close()

	chUp := make(chan int64)
	chDown := make(chan int64)
	go func() {
		// copy from conn to upstream
		n, err := io.Copy(upstreamConn, conn)
		if err != nil {
			// TODO: log error
			_ = err
		}
		chUp <- n
	}()
	go func() {
		// copy from upstream to conn
		n, err := io.Copy(conn, upstreamConn)
		if err != nil {
			// TODO: log error
			_ = err
		}
		chDown <- n
	}()

	// Collect statistics
	uploadBytes := <-chUp
	downloadBytes := <-chDown
	// TODO: print or log them?
	_, _ = uploadBytes, downloadBytes
}

func (s *HTTPServer) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.Handle(conn.(*net.TCPConn))
	}
}

func (s *HTTPServer) ListenAndServe() error {
	l, err := net.Listen("tcp", s.s.httpAddr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}
