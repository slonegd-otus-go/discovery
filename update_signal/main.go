package main



import (
	"log"
	"net/http"
	"time"

	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/slonegd-otus-go/discovery/registry/nats"

	"github.com/slonegd-otus-go/discovery"
)

func main()  {
	registry := nats.NewRegistry(registry.Addrs("nats://127.0.0.1:4222")) // TODO вынести в глобалку/флаги
	selector := selector.NewSelector(selector.Registry(registry))
	discovery := discovery.New(selector)
	address, err := discovery.GetAddress("sedmax.web.signals")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("address of sedmax.web.signals:", address)

	response, err := http.Get("http://"+address + "/update/configuration")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("send update configuration signal, response: %+v", response)
	time.Sleep(time.Second)
}
