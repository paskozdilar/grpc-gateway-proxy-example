package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jnis23/grpc-gateway-proxy-example/stream"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGrpcServer() {
	server := grpc.NewServer()
	stream.RegisterStreamServiceServer(server, &StreamServer{})
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf("grpc server failed : %v", err)
	}
}

func runHttpGateway() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := stream.RegisterStreamServiceHandlerFromEndpoint(context.Background(), mux, "localhost:9090", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	mux.HandlePath("GET", "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.ServeFile(w, r, "index.html")
	})

	ws := wsproxy.WebsocketProxy(mux)
	if err := http.ListenAndServe(":8080", ws); err != nil {
		log.Fatalf("gateway server failed: %v", err)
	}
}

func main() {
	go runGrpcServer()
	go runHttpGateway()
	fmt.Println("servers started")
	select {}
}
