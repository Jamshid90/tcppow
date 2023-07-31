package pow

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrMaxExceeded = errors.New("max iterations exceeded")
	ErrParseing    = errors.New("error parseing hashdata")
)

const zeroByte = 48

type HashData struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
}

func (h *HashData) ComputeHash(maxIterations int) error {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		header := h.String()
		hash := sha1Hash(header)
		if IsHashCorrect(hash, h.ZerosCount) {
			return nil
		}
		h.Counter++
	}
	return ErrMaxExceeded
}

func (h HashData) String() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

func Parse(hashDataStr []byte) (*HashData, error) {
	var hashData HashData
	err := json.Unmarshal([]byte(hashDataStr), &hashData)
	if err != nil {
		return nil, ErrParseing
	}
	return &hashData, nil
}
