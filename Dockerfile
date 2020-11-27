FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go build -v cmd/server/server.go 
# RUN go install -v ./...
RUN go run cmd/csv2protobuf/csv_to_protobuf.go -o assets/dervaze-rootset.protobuf -i assets/rootdata/
EXPOSE 9876
CMD ["./server", "-i", "assets/dervaze-rootset.protobuf", "-h", "0.0.0.0", "-p", "9876"]
# CMD ["/bin/bash"]

