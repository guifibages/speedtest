#!/usr/bin/env bash

source $GOROOT/src/golang-crosscompile/crosscompile.bash
mkdir build
for arch in darwin-amd64 linux-amd64 linux-arm linux-386
do
	go-${arch} build -o build/speedtest-${arch} st.go
done
for arch in windows-386 windows-amd64
do
	go-${arch} build -o build/speedtest-${arch}.exe st.go
done
