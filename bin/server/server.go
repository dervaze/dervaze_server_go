package main

import (
	dervaze "dervaze/lang"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var (
	inputfile  = flag.String("i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	host       = flag.String("h", "0.0.0.0", "IP address or hostname to listen to")
	port       = flag.Int("p", 9876, "port to listen to")
	serverType = flag.String("s", "REST", "Server type to start. Can be REST or GRPC. You can use $SERVER_TYPE environment variable argument to set this as well.")
)

func main() {

	if serverTypeEnv, present := os.LookupEnv("SERVER_TYPE"); present {
		flag.Set("s", serverTypeEnv)
	}

	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s=%s [%s] \n", f.Name, f.Value.String(), f.Usage)
	})

	dervaze.InitSearch(*inputfile)
	*serverType = strings.ToLower(*serverType)
	if *serverType == "grpc" {
		fmt.Println("Starting GRPC Server")
		grpcServer(*host, *port)
	} else if *serverType == "rest" || *serverType == "json" {
		restServer(*host, *port)
		fmt.Println("Starting REST Server")
	} else {
		panic("Need either REST or GRPC as server type: " + *serverType)
	}
}

func restServer(host string, port int) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/json/prefix/tr/{word}", dervaze.JSONPrefixTr)
	router.HandleFunc("/v1/json/prefix/ot/{word}", dervaze.JSONPrefixOt)
	router.HandleFunc("/v1/json/exact/tr/{word}", dervaze.JSONExactTr)
	router.HandleFunc("/v1/json/exact/ot/{word}", dervaze.JSONExactOt)
	router.HandleFunc("/v1/json/search/any/{word}", dervaze.JSONSearchAuto)
	router.HandleFunc("/v1/json/search/ot/{word}", dervaze.JSONSearchOt)
	router.HandleFunc("/v1/json/search/tr/{word}", dervaze.JSONSearchTr)
	router.HandleFunc("/v1/json/exact/abjad/{number}", dervaze.JSONExactAbjad)
	router.HandleFunc("/v1/json/calc/abjad/{word}", dervaze.JSONCalcAbjad)
	router.HandleFunc("/v1/json/v2u/{word}", dervaze.JSONV2U)
	router.HandleFunc("/v1/json/u2v/{word}", dervaze.JSONU2V)
	router.HandleFunc("/v1/version/", dervaze.JSONVersion)
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

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
