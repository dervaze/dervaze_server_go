package main

import (
	dervaze "dervaze/lang"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	inputfile = flag.String("i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	host      = flag.String("h", "127.0.0.1", "IP address or hostname to listen to")
	port      = flag.Int("p", 9876, "port to listen to")
)

func server(host string, port int) {

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

func main() {

	flag.Parse()

	fmt.Println("Starting gRPC Server")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s=%s [%s] \n", f.Name, f.Value.String(), f.Usage)
	})

	dervaze.InitSearch(*inputfile)
	server(*host, *port)
}
