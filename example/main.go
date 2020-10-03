package main

import (
	"log"
	"os"

	"github.com/slonegd-otus-go/discovery"
	"github.com/slonegd-otus-go/discovery/registry/nats"

	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

func main() {
	service := "sedmax.controller"
	if len(os.Args) > 1 {
		service = os.Args[1]
	}
	registry := nats.NewRegistry(registry.Addrs("nats://127.0.0.1:4222"))
	selector := selector.NewSelector(selector.Registry(registry))
	discovery := discovery.New(selector)

	address, err := discovery.GetAddress(service)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(address)
}
