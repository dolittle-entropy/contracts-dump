package main

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MessageOrigin string

const (
	HeadMessage    MessageOrigin = "head"
	RuntimeMessage MessageOrigin = "runtime"
)

type ContractsStream struct {
	Mutex          sync.Mutex
	OriginalStream grpc.ServerStream
	Parser         *Parser
	CurrentCall    uint64
	FullMethod     string
}

func (cs *ContractsStream) Context() context.Context {
	return cs.OriginalStream.Context()
}

func (cs *ContractsStream) SetHeader(headers metadata.MD) error {
	return cs.OriginalStream.SetHeader(headers)
}

func (cs *ContractsStream) SendHeader(headers metadata.MD) error {
	return cs.OriginalStream.SendHeader(headers)
}

func (cs *ContractsStream) SetTrailer(trailers metadata.MD) {
	cs.OriginalStream.SetTrailer(trailers)
}

func (cs *ContractsStream) SendMsg(m interface{}) error {
	if message, err := cs.Parser.ParseMessage(cs.FullMethod, RuntimeMessage, m); err == nil {
		dump(Message{
			CallCounter: cs.CurrentCall,
			Timestamp:   time.Now(),
			Origin:      RuntimeMessage,
			Message:     message,
		}, cs.CurrentCall)
	} else {
		logrus.WithError(err).Warn("Failed to decode sent message")
	}

	return cs.OriginalStream.SendMsg(m)
}

func (cs *ContractsStream) RecvMsg(m interface{}) error {
	err := cs.OriginalStream.RecvMsg(m)
	if err != nil {
		return err
	}

	if message, err := cs.Parser.ParseMessage(cs.FullMethod, HeadMessage, m); err == nil {
		dump(Message{
			CallCounter: cs.CurrentCall,
			Timestamp:   time.Now(),
			Origin:      HeadMessage,
			Message:     message,
		}, cs.CurrentCall)
	} else {
		logrus.WithError(err).Warn("Failed to decode received message")
	}

	return nil
}
