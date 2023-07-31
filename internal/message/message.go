package message

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	Quit = iota
	RequestChallenge
	ResponseChallenge
	RequestResource
	ResponseResource
)

var (
	ErrInvalidMsg     = errors.New("message is invalid")
	ErrInvalidMsgType = errors.New("message type is invalid")
	ErrWrite          = errors.New("error writeing message")
	ErrRead           = errors.New("error reading message")
)

type Message struct {
	Type int
	Body string
}

func (m *Message) String() string {
	return fmt.Sprintf("%d|%s", m.Type, m.Body)
}

func Parse(str string) (*Message, error) {
	var (
		typeIndex = 0
		bodyIndex = 1
	)
	msgSlices := strings.Split(strings.TrimSpace(str), "|")
	if len(msgSlices) < 1 || len(msgSlices) > 2 {
		return nil, ErrInvalidMsg
	}

	mType, err := strconv.Atoi(msgSlices[typeIndex])
	if err != nil {
		return nil, ErrInvalidMsgType
	}
	msg := Message{
		Type: mType,
	}

	if len(msgSlices) == 2 {
		msg.Body = msgSlices[bodyIndex]
	}

	return &msg, nil
}

func Send(msg *Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.String())
	_, err := conn.Write([]byte(msgStr))
	if err != nil {
		return ErrWrite
	}
	return nil
}

func Read(reader *bufio.Reader) (*Message, error) {
	msgStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, ErrRead
	}
	msg, err := Parse(msgStr)
	if err != nil {
		return nil, fmt.Errorf("err parse msg: %w", err)
	}
	return msg, nil
}
