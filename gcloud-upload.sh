#!/bin/bash

gcloud builds submit --tag gcr.io/dervaze-server-go/v0.1
gcloud run deploy --image gcr.io/dervaze-server-go/v0.1 --platform managed --port 9876 --memory 512M
