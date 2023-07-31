package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Jamshid90/tcppow/internal/message"
)

type Handler interface {
	HandleConnection(ctx context.Context, conn net.Conn)
	RequestChallenge(ctx context.Context) (*message.Message, error)
	RequestResource(ctx context.Context, msg *message.Message) (*message.Message, error)
}

type Server struct {
	address string
	handler Handler
}

func New(host, port string, handler Handler) *Server {
	return &Server{
		address: fmt.Sprintf("%s:%s", host, port),
		handler: handler,
	}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			fmt.Println("error connection accept: %w", err)
			continue
		}
		go s.handler.HandleConnection(ctx, conn)
	}
}
