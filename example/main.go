package main

import (
	"log"
	"os"

	"git.sedmax.ru/CORE/sed_controller/tests/discovery"

	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/nats/v2"
)

func main() {
	service := "sedmax.controller"
	if len(os.Args) > 1 {
		service = os.Args[1]
	}
	registry := nats.NewRegistry(registry.Addrs("nats://127.0.0.1:4222")) // TODO вынести в глобалку/флаги
	selector := selector.NewSelector(selector.Registry(registry))
	discovery := discovery.New(selector)

	address, err := discovery.GetAddress(service)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(address)
}
