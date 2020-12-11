protoc -I . --grpc-gateway_out . \                                                                                   ‹master⚡› ➤ fc17f69  [11m] 
     --grpc-gateway_opt logtostderr=true \
     --grpc-gateway_opt paths=source_relative \
     --grpc-gateway_opt generate_unbound_methods=true \
     lang/dervaze.proto

     protoc --go_out=lang --go_opt=paths=source_relative  --go-grpc_out=lang --go-grpc_opt=paths=source_relative lang/dervaze.proto