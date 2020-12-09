package main

import (
	"context"
	dervaze "dervaze/lang"
	"flag"
	"fmt"
	"net"
	"regexp"
	"strconv"

	"google.golang.org/grpc"
)

var (
	inputfile = flag.String("i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	host      = flag.String("h", "127.0.0.1", "IP address or hostname to listen to")
	port      = flag.Int("p", "9876", "port to listen to")
)

type dervazeServer struct {
}

func (*dervazeServer) VisencToOttoman(ctx context.Context, in *dervaze.OttomanWord, opts ...grpc.CallOption) (*dervaze.OttomanWord, error) {

	out := dervaze.MakeOttomanWord(in.Visenc, "")
	return &out, nil

}
func (*dervazeServer) OttomanToVisenc(ctx context.Context, in *dervaze.OttomanWord, opts ...grpc.CallOption) (*dervaze.OttomanWord, error) {
	out := dervaze.MakeOttomanWord("", in.Unicode)
	return &out, nil
}
func (*dervazeServer) SearchRoots(ctx context.Context, in *dervaze.SearchRequest, opts ...grpc.CallOption) (*dervaze.RootSet, error) {

	var rootList []*dervaze.Root
	var err error = nil

	searchField := in.SearchField
	searchString := in.SearchString
	maxLen := int(in.ResultLimit)

	switch in.SearchType {
	case dervaze.SearchType_FUZZY:
		switch searchField {
		case dervaze.SearchField_AUTO:
			rootList = dervaze.FuzzySearchAuto(searchString, maxLen)
		case dervaze.SearchField_OTTOMAN:
			rootList = dervaze.FuzzySearchUnicode(searchString, maxLen)
		case dervaze.SearchField_TURKISH_LATIN:
			rootList = dervaze.FuzzySearchTurkishLatin(searchString, maxLen)
		case dervaze.SearchField_VISENC:
			rootList = dervaze.FuzzySearchVisenc(searchString, maxLen)
		case dervaze.SearchField_ABJAD:
			if s, e := strconv.Atoi(searchString); e == nil {
				rootList = dervaze.IndexSearchAbjad(int32(s), maxLen)
			} else {
				err = e
			}
		}
	case dervaze.SearchType_REGEX:
		if searchRegex, e := regexp.Compile(searchString); e == nil {

			switch searchField {
			case dervaze.SearchField_AUTO:
				rootList = dervaze.RegexSearchAuto(searchRegex, maxLen)
			case dervaze.SearchField_OTTOMAN:
				rootList = dervaze.RegexSearchUnicode(searchRegex, maxLen)
			case dervaze.SearchField_TURKISH_LATIN:
				rootList = dervaze.RegexSearchTurkishLatin(searchRegex, maxLen)
			case dervaze.SearchField_VISENC:
				rootList = dervaze.RegexSearchVisenc(searchRegex, maxLen)
			case dervaze.SearchField_ABJAD:
				if s, e := strconv.Atoi(searchString); e == nil {
					rootList = dervaze.IndexSearchAbjad(int32(s), maxLen)
				} else {
					err = e
				}
			}
		} else {
			err = e
		}

	case dervaze.SearchType_PREFIX:

		switch searchField {
		case dervaze.SearchField_AUTO:
			rootList = dervaze.PrefixSearchAuto(searchString, maxLen)
		case dervaze.SearchField_OTTOMAN:
			rootList = dervaze.PrefixSearchUnicode(searchString, maxLen)
		case dervaze.SearchField_TURKISH_LATIN:
			rootList = dervaze.PrefixSearchTurkishLatin(searchString, maxLen)
		case dervaze.SearchField_VISENC:
			rootList = dervaze.PrefixSearchVisenc(searchString, maxLen)
		case dervaze.SearchField_ABJAD:
			if s, e := strconv.Atoi(searchString); e == nil {
				rootList = dervaze.IndexSearchAbjad(int32(s), maxLen)
			} else {
				err = e
			}
		}

	}

}
func (*dervazeServer) Translate(ctx context.Context, in *dervaze.TranslateRequest, opts ...grpc.CallOption) (*dervaze.TranslateResponse, error) {

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

	flag.Parse()

	fmt.Println("Starting gRPC Server")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s=%s [%s] \n", f.Name, f.Value.String(), f.Usage)
	})

	dervaze.InitSearch(inputfile)
	server(host, port)
}
