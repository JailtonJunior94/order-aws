package entities

import (
	"time"
)

type (
	Sequence struct {
		Date     string
		Sequence int64
		ExpireAt int64
		Code     string
	}
)

func NewSequence(code string) *Sequence {
	return &Sequence{
		Sequence: 1,
		Code:     code,
		Date:     time.Now().Format("20060102"),
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	}
}
