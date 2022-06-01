package main

import (
	"fmt"
	"grpc-coba/data"
	"grpc-coba/server"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	protos "grpc-coba/protos/currency"
)

func main() {
	log := hclog.Default()
	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("Unnable to generate rates", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	cs := server.NewCurrency(rates, log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9091))
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
