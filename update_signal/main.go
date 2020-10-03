package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	http_client "github.com/micro/go-plugins/client/http/v2"

	"github.com/slonegd-otus-go/discovery"
	"github.com/slonegd-otus-go/discovery/registry/nats"
)

func main() {
	registry := nats.NewRegistry(registry.Addrs("nats://127.0.0.1:4222")) // TODO вынести в глобалку/флаги
	selector := selector.NewSelector(selector.Registry(registry))
	discovery := discovery.New(selector)
	address, err := discovery.GetAddress("sedmax.web.signals")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("address of sedmax.web.signals:", address)

	log.Printf("send update configuration signal via http")

	response, err := http.Get("http://" + address + "/update/configuration")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response: %+v", response)

	log.Printf("send update configuration signal via micro")

	service := micro.NewService(
		micro.Name("test.update_signal"),
		micro.Registry(registry),
		micro.Client(http_client.NewClient()),
	)

	request := service.Client().NewRequest("sedmax.web.signals", "/update/configuration", nil)
	var microResponse interface{}
	err = service.Client().Call(context.Background(), request, &microResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response: %+v", microResponse)

	time.Sleep(time.Second)
}
