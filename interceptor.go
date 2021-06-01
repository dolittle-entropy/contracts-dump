package main

import (
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

func ContractsInterceptor(parser *Parser) grpc.StreamServerInterceptor {
	var callCounter uint64

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		currentCall := atomic.AddUint64(&callCounter, 1)

		fullMethod := strings.Split(info.FullMethod, "/")
		dump(Call{
			CallCounter: currentCall,
			Timestamp:   time.Now(),
			Service:     fullMethod[1],
			Method:      fullMethod[2],
		}, currentCall)

		stream := &ContractsStream{
			OriginalStream: ss,
			Parser:         parser,
			CurrentCall:    currentCall,
			FullMethod:     info.FullMethod,
		}

		err := handler(srv, stream)
		if err != nil {
			dump(Result{
				CallCounter: currentCall,
				Timestamp:   time.Now(),
				Success:     false,
				Error:       err,
			}, currentCall)
			return err
		}

		dump(Result{
			CallCounter: currentCall,
			Timestamp:   time.Now(),
			Success:     true,
		}, currentCall)
		return nil
	}
}
