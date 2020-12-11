FROM golang:1.15

WORKDIR /go/src/app
COPY . .

RUN go build -v bin/rest_server/rest_server.go 
RUN go build -v bin/csv2protobuf/csv_to_protobuf.go
# RUN go install -v ./...
RUN ./csv_to_protobuf -i assets/rootdata/ -r assets/dervaze-rootset.protobuf -s assets/dervaze-suffixset.protobuf -f protobuf
EXPOSE 9876
RUN ls -R >> /dev/log
ENTRYPOINT ["./rest_server", "-i", "assets/dervaze-rootset.protobuf", "-h", "0.0.0.0", "-p", "9876"]
# CMD ["/bin/bash"]

