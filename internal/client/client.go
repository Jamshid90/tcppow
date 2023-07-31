package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/Jamshid90/tcppow/internal/message"
	"github.com/Jamshid90/tcppow/internal/pkg/pow"
)

type Client struct {
	address string
}

func New(host, port string) *Client {
	return &Client{
		address: fmt.Sprintf("%s:%s", host, port),
	}
}

func (c *Client) Run(ctx context.Context, hashMaxIterations int) error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("error dial: %w", err)
	}

	defer conn.Close()

	for {
		message, err := HandleConnection(ctx, conn, conn, hashMaxIterations)
		if err != nil {
			return err
		}
		fmt.Println("quote result:", message)
	}
}

func HandleConnection(ctx context.Context, readerConn io.Reader, writerConn io.Writer, hashMaxIterations int) (string, error) {

	// 1. request challenge
	err := message.Send(&message.Message{
		Type: message.RequestChallenge,
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("error send request challenge: %w", err)
	}

	// reading and parsing response
	reader := bufio.NewReader(readerConn)
	msg, err := message.Read(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	fmt.Println("got message", msg)

	// 2. got challenge, compute hashdata
	hashdata, err := pow.Parse([]byte(msg.Body))
	if err != nil {
		return "", fmt.Errorf("err parse hashdata: %w", err)
	}
	if err = hashdata.ComputeHash(hashMaxIterations); err != nil {
		return "", fmt.Errorf("err compute hashdata: %w", err)
	}
	byteData, err := json.Marshal(hashdata)
	if err != nil {
		return "", fmt.Errorf("err marshal hashdata: %w", err)
	}

	fmt.Println("challenge sent to server", hashdata)

	// 3. send challenge solution back to server
	err = message.Send(&message.Message{
		Type: message.RequestResource,
		Body: string(byteData),
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	// 4. get result quote from server
	msg, err = message.Read(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	return msg.Body, nil
}
