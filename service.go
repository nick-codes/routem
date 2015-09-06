package routem

import (
	"crypto/tls"
	"net"
	"net/http"
)

type (
	service struct {
		address  string
		listener net.Listener
		server   *http.Server
		err      error
		started  chan struct{}
		running  chan struct{}
	}
)

func (s *service) Address() string {
	return s.address
}

func (s *service) IsRunning() bool {
	select {
	case <-s.running:
		return false
	default:
		return true
	}
}

func (s *service) Wait() error {
	<-s.running
	return s.err
}

func (s *service) Stop() error {
	return s.listener.Close()
}

func (s *service) run() error {
	listener, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	s.listener = listener

	s.serve()

	return nil
}

func (s *service) runTLS(certFile, keyFile string) error {
	config := &tls.Config{}

	config.NextProtos = []string{"http/1.1"}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	s.listener = tls.NewListener(listener, config)

	s.serve()

	return nil
}

func (s *service) serve() {
	go func() {
		close(s.started)
		s.err = s.server.Serve(s.listener)
		close(s.running)
	}()
	<-s.started
}

func newService(address string, handler http.Handler) *service {
	s := &service{
		address: address,
		server: &http.Server{
			Addr:    address,
			Handler: handler,
		},
		started: make(chan struct{}),
		running: make(chan struct{}),
	}

	return s
}
