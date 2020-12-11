package main

import (
	dervaze "dervaze/lang"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func grpcServer(host string, port int) {

	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	dd := dervaze.NewDervazeServerImpl()

	dervaze.RegisterDervazeServer(server, dd)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
