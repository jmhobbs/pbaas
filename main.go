package main

import (
	"encoding/base64"
	"encoding/binary"
	"flag"
	"log"
	"math/rand"
)

var (
	grpcAddress string
	httpAddress string
)

func init() {
	flag.StringVar(&grpcAddress, "grpc", ":12345", "Address for gRPC to listen on")
	flag.StringVar(&httpAddress, "http", ":9999", "Address for HTTP to listen on")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 {
		flag.Usage()
		return
	}

	store := NewInMemoryProgressDB()

	grpc := progressBarServiceServer{store}
	log.Printf("Serving gRPC on %s\n", grpcAddress)
	go grpc.Serve(grpcAddress)

	log.Printf("Serving HTTP on %s\n", httpAddress)
	http := NewWebServer(store)
	go http.Serve(httpAddress)

	done := make(chan bool)
	<-done
}

func newToken() string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(b)
}
