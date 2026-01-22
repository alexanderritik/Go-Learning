package ingestion

import (
	"bufio"
	"log"
	"net"
	"sync"
)

type Server struct {
	listenAddr string

	ln net.Listener

	DataCh chan []byte

	quitCh chan struct{}

	wg sync.WaitGroup
}

func NewServer(addr string, bufferSize int) *Server {
	return &Server{
		listenAddr: addr,
		DataCh:     make(chan []byte, bufferSize),
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Println("Error starting server:", err)
		return err
	}
	s.ln = ln

	go s.loop()

	log.Printf("LogHive Ingestion Server running on %s", s.listenAddr)
	return nil
}

func (s *Server) loop() {
	defer s.ln.Close()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.quitCh:
				return
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}
		s.wg.Add(1)

		go s.handleConn(conn)

	}
}
func (s *Server) handleConn(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Bytes()
		msg := make([]byte, len(text))
		copy(msg, text)

		select {
		case s.DataCh <- msg:
			// Success
		case <-s.quitCh:
			// Server is shutting down, stop reading
			return
		}
	}
}

func (s *Server) Stop() {
	close(s.quitCh) // Signal all components to stop
	s.ln.Close()    // Stop accepting new connections
	s.wg.Wait()     // Wait for existing connections to finish reading
	close(s.DataCh) // Close the data channel so consumers know no more data is coming
	log.Println("Ingestion Server stopped")
}
