#!/bin/bash

cp bin/grpc_server/Dockerfile .
gcloud builds submit --tag gcr.io/dervaze-grpc-go/v0.1
gcloud run deploy --image gcr.io/dervaze-grpc-go/v0.1 --memory 512M --platform managed --port 9876 

