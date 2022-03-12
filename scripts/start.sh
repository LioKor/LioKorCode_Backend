#!/bin/bash

mkdir build/
go build -o build/main_service cmd/main.go

nohup ./build/main_service > main.out &
