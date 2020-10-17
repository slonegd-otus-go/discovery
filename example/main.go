package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	http_client "github.com/micro/go-plugins/client/http/v2"

	"github.com/slonegd-otus-go/discovery/discovery"
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

	log.Printf("read archive via micro")

	// в плагине http клиента есть фильтр, который ищет только http сервера
	// соответсвенно сервера должны быть как то помечены, но у нас нет
	// поэтому всегда не находит, но если его закомментировать, то работает
	// only get the things that are of mucp protocol
	// selectOptions := append(opts.SelectOptions, selector.WithFilter(
	// 	selector.FilterLabel("protocol", "http"),
	// ))
	client := http_client.NewClient(
		client.Selector(selector),
		client.ContentType("application/json"),
	)
	request := client.NewRequest("sedmax.web.signals", "/read_archive", "{}")

	var microResponse interface{}
	err = client.Call(context.Background(), request, &microResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response: %+v", microResponse)

	time.Sleep(time.Second)
}

type ReadArchiveCommand struct {
	// Идентификатор устройства
	// Required: true
	// Example: 1000
	DeviceID int `json:"device_id"`
	// Начало интервала, дата в формате YYYY-MM-DD HH:MI:SS
	// Required: true
	// Example: 2019-03-21
	BeginDateTime string `json:"begin"`
	// Конец интервала, дата в формате YYYY-MM-DD HH:MI:SS
	// Required: true
	// Example: 2019-03-22
	EndDateTime string `json:"end"`
	// Читать или не читать почасовые архивы
	// Required: true
	IsReadArchive bool `json:"is_read_archive"`
	// Читать или не читать журналы событий
	// Required: true
	IsReadEventJournal bool `json:"is_read_event_journal"`
}
