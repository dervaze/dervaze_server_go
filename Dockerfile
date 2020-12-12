FROM golang:1.15

ARG SERVER_TYPE
ENV SERVER_TYPE=${SERVER_TYPE}
WORKDIR /go/src/app
COPY . .
RUN go build -v bin/server/server.go 
RUN go build -v bin/csv2protobuf/csv_to_protobuf.go
# RUN go install -v ./...
RUN ./csv_to_protobuf -i assets/rootdata/ -r assets/dervaze-rootset.protobuf -s assets/dervaze-suffixset.protobuf -f protobuf
EXPOSE 9876
# RUN ls -R 
RUN env
ENTRYPOINT ["./server", "-i", "assets/dervaze-rootset.protobuf", "-h", "0.0.0.0", "-p", "9876"]
# CMD ["/bin/bash"]

