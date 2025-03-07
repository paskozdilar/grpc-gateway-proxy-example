package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jnis23/grpc-gateway-proxy-example/stream"
)

type StreamServer struct {
	stream.UnimplementedStreamServiceServer
}

func (s *StreamServer) Stream(req *stream.StreamRequest, str stream.StreamService_StreamServer) error {
	ctx := str.Context()
	fmt.Printf("stream started for %d\n", req.Id)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done!")
			return ctx.Err()
		case <-ticker.C:
			msg := &stream.StreamResponse{
				Message: "tick",
			}
			if err := str.Send(msg); err != nil {
				return err
			}
		}
	}
}

func (s *StreamServer) EmptyStream(req *empty.Empty, str stream.StreamService_EmptyStreamServer) error {
	ctx := str.Context()
	fmt.Println("empty stream started")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done!")
			return ctx.Err()
		case <-ticker.C:
			msg := &stream.StreamResponse{
				Message: "tick",
			}
			if err := str.Send(msg); err != nil {
				return err
			}
		}
	}
}

func (s *StreamServer) Blocker(req *empty.Empty, str stream.StreamService_BlockerServer) error {
	ctx := str.Context()
	fmt.Println("blocker stream started")
	<-ctx.Done()
	fmt.Println("context done!")
	return nil
}
