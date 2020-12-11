package main

import (
	dervaze "dervaze/lang"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	gmux "github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	inputfile = flag.String("i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	host      = flag.String("h", "0.0.0.0", "IP address or hostname to listen to")
	port      = flag.Int("p", 9876, "port to listen to")
)

type myService struct{}

func (m *myService) Echo(c context.Context, s *pb.EchoMessage) (*pb.EchoMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func commonServer(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	// opts := []grpc.ServerOption{
	// 	grpc.Creds(credentials.NewClientTLSFromCert(demoCertPool, addr))}
	//
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	dervaze.RegisterDervazeServer(grpcServer, dervaze.NewDervazeServerImpl())
	ctx := context.Background()

	// dcreds := credentials.NewTLS(&tls.Config{
	// 	ServerName: demoAddr,
	// 	RootCAs:    demoCertPool,
	// })
	// dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	//
	dopts := []grpc.DialOption{}

	router := gmux.NewRouter().StrictSlash(true)
	// mux := http.NewServeMux()
	router.HandleFunc("/v1/json/prefix/tr/{word}", JSONPrefixTr)
	router.HandleFunc("/v1/json/prefix/ot/{word}", JSONPrefixOt)
	router.HandleFunc("/v1/json/exact/tr/{word}", JSONExactTr)
	router.HandleFunc("/v1/json/exact/ot/{word}", JSONExactOt)
	router.HandleFunc("/v1/json/search/any/{word}", JSONSearchAuto)
	router.HandleFunc("/v1/json/search/ot/{word}", JSONSearchOt)
	router.HandleFunc("/v1/json/search/tr/{word}", JSONSearchTr)
	router.HandleFunc("/v1/json/exact/abjad/{number}", JSONExactAbjad)
	router.HandleFunc("/v1/json/calc/abjad/{word}", JSONCalcAbjad)
	router.HandleFunc("/v1/json/v2u/{word}", JSONV2U)
	router.HandleFunc("/v1/json/u2v/{word}", JSONU2V)
	router.HandleFunc("/v1/version/", JSONVersion)

	gwmux := runtime.NewServeMux()
	err := dervaze.RegisterDervazeHandlerFromEndpoint(ctx, gwmux, addr, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	router.Handle("/", gwmux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, router),
		// TLSConfig: &tls.Config{
		// 	Certificates: []tls.Certificate{*demoKeyPair},
		// 	NextProtos:   []string{"h2"},
		// },
	}

	fmt.Printf("grpc on port: %d\n", port)
	// err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	err = srv.Serve(conn)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}

func main() {

	flag.Parse()

	fmt.Println("Starting Common (gRPC + REST) Server")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s=%s [%s] \n", f.Name, f.Value.String(), f.Usage)
	})

	dervaze.InitSearch(*inputfile)
	server(*host, *port)
}
