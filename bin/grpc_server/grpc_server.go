package main

import (
	"context"
	dervaze "dervaze/lang"
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

var (
	
	inputfile = flag.String("i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	flag.StringVar(&host, "h", "127.0.0.1", "IP address or hostname to listen to")
	flag.StringVar(&port, "p", "9876", "port to listen to")
)

type dervazeServer struct {
}

func (*dervazeServer) VisencToOttoman(ctx context.Context, in *OttomanWord, opts ...grpc.CallOption) (*OttomanWord, error) {

}
func (*dervazeServer) OttomanToVisenc(ctx context.Context, in *OttomanWord, opts ...grpc.CallOption) (*OttomanWord, error) {

}
func (*dervazeServer) SearchRoots(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*RootSet, error) {

}
func (*dervazeServer) Translate(ctx context.Context, in *TranslateRequest, opts ...grpc.CallOption) (*TranslateResponse, error) {

}
func server(host, port string) {

	listener, err := net.Listen("tcp", ":9877")

	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	ds := dervazeServer{}

	dervaze.RegisterDervazeServer(server, &ds)

}

func main() {

	var inputfile string
	var port string
	var host string

	flag.StringVar(&inputfile, "i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	flag.StringVar(&host, "h", "127.0.0.1", "IP address or hostname to listen to")
	flag.StringVar(&port, "p", "9876", "port to listen to")

	flag.Parse()

	fmt.Println("Starting gRPC Server")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s=%s [%s] \n", f.Name, f.Value.String(), f.Usage)
	})

	dervaze.InitSearch(inputfile)
	server(host, port)
}
