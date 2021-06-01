package main

import (
	"time"

	"github.com/jhump/protoreflect/dynamic"
)

type Call struct {
	CallCounter uint64    `json:"call"`
	Timestamp   time.Time `json:"timestamp"`
	Service     string    `json:"service"`
	Method      string    `json:"method"`
}

type Message struct {
	CallCounter uint64           `json:"call"`
	Timestamp   time.Time        `json:"timestamp"`
	Origin      MessageOrigin    `json:"origin"`
	Message     *dynamic.Message `json:"message"`
}

type Result struct {
	CallCounter uint64    `json:"call"`
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
	Error       error     `json:"error,omitempty"`
}
