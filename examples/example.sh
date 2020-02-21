#!/bin/bash

go build -o myexample examples/example.go
./myexample -peer "localhost:8081,localhost:8082" -port 8080 &
./myexample -peer "localhost:8080,localhost:8082" -port 8081 &
./myexample -peer "localhost:8081,localhost:8080" -port 8082
