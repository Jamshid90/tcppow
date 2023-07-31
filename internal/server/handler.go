package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/Jamshid90/tcppow/internal/message"
	"github.com/Jamshid90/tcppow/internal/pkg/pow"
)

type ctxKeyClientInfo int

const (
	keyClientInfo ctxKeyClientInfo = 0
)

var quotes = []string{
	"Science is organized knowledge. Wisdom is organized life.",
	"Doubt is the origin of wisdom.",
	"The truest wisdom is a resolute determination.",
	"Wisdom is not a product of schooling but of the lifelong attempt to acquire it.",
	"Wisdom is the power to put our time and our knowledge to the proper use.",
	"A symptom of wisdom is curiosity. The evidence is calmness and perseverance. The causes are experimentation and understanding.",
	"It is not the man who has too little, but the man who craves more, that is poor.",
	"A wise man never loses anything, if he has himself.",
	"A fool is known by his speech; and a wise man by silence.",
	"Foolishness is a twin sister of wisdom. ",
	"There is only a finger’s difference between a wise man and a fool.",
	"Never say no twice if you mean it. ",
}

type Cache interface {
	Set(ctx context.Context, key string, expiration time.Duration) error
	Get(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
}

type handler struct {
	cache          Cache
	hashDuration   time.Duration
	hashZerosCount int
}

func NewHandler(cache Cache, hashDuration time.Duration, hashZerosCount int) *handler {
	return &handler{
		cache:          cache,
		hashDuration:   hashDuration,
		hashZerosCount: hashZerosCount,
	}
}

func (h *handler) HandleConnection(ctx context.Context, conn net.Conn) {

	fmt.Println("new connection:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		req, err := message.Read(reader)
		if err != nil {
			fmt.Println("err read msg:", err)
			return
		}

		var msg *message.Message
		ctx = context.WithValue(ctx, keyClientInfo, conn.RemoteAddr().String())

		switch req.Type {
		case message.RequestChallenge:
			msg, err = h.RequestChallenge(ctx)
			if err != nil {
				fmt.Println("error request challenge:", err)
				return
			}
		case message.RequestResource:
			msg, err = h.RequestResource(ctx, req)
			if err != nil {
				fmt.Println("error request resource:", err)
				return
			}
		default:
			fmt.Println("error unknown type message")
			return
		}

		if msg != nil {
			if err = message.Send(msg, conn); err != nil {
				fmt.Println("error send message:", err)
			}
		}
	}
}

func (h *handler) RequestChallenge(ctx context.Context) (*message.Message, error) {
	randValue := strconv.Itoa(rand.Intn(100000))
	err := h.cache.Set(ctx, randValue, h.hashDuration)
	if err != nil {
		return nil, fmt.Errorf("error add rand to cache: %w", err)
	}

	clientInfo, ok := ctx.Value(keyClientInfo).(string)
	if !ok {
		return nil, fmt.Errorf("error client information getting")
	}

	hashdata := pow.HashData{
		Version:    1,
		ZerosCount: h.hashZerosCount,
		Date:       time.Now().Unix(),
		Resource:   clientInfo,
		Rand:       base64.StdEncoding.EncodeToString([]byte(randValue)),
		Counter:    0,
	}

	hashJB, err := json.Marshal(hashdata)
	if err != nil {
		return nil, fmt.Errorf("error marshal hash: %v", err)
	}

	msg := message.Message{
		Type: message.ResponseChallenge,
		Body: string(hashJB),
	}

	return &msg, nil
}

func (h *handler) RequestResource(ctx context.Context, msg *message.Message) (*message.Message, error) {

	var hashdata pow.HashData
	err := json.Unmarshal([]byte(msg.Body), &hashdata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal hash: %w", err)
	}

	clientInfo, ok := ctx.Value(keyClientInfo).(string)
	if !ok {
		return nil, fmt.Errorf("error client information getting")
	}

	if hashdata.Resource != clientInfo {
		return nil, fmt.Errorf("invalid hash resource")
	}

	randValue, err := base64.StdEncoding.DecodeString(hashdata.Rand)
	if err != nil {
		return nil, fmt.Errorf("error decode rand: %w", err)
	}

	exists, err := h.cache.Get(ctx, string(randValue))
	if err != nil {
		return nil, fmt.Errorf("error get rand from cache: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("error challenge expired or not sent")
	}

	if time.Duration(time.Now().Unix()-hashdata.Date) > h.hashDuration {
		return nil, fmt.Errorf("error challenge expired")
	}

	maxIter := hashdata.Counter
	if maxIter == 0 {
		maxIter = 1
	}

	err = hashdata.ComputeHash(maxIter)
	if err != nil {
		return nil, fmt.Errorf("invalid hash")
	}

	return &message.Message{
		Type: message.ResponseResource,
		Body: quotes[rand.Intn(11)],
	}, nil
}
