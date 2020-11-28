#!/bin/bash

sudo docker build . -t dervaze_server_go
sudo docker run -p 9876:9876  -t dervaze_server_go
