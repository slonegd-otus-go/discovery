package json

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/micro/go-micro/v2/registry"
)

type Service struct {
	Name      string               `json:"name"`
	Version   string               `json:"version"`
	Metadata  map[string]string    `json:"metadata"`
	Endpoints []*registry.Endpoint `json:"endpoints"`
	Nodes     []*Node              `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
}

type Result struct {
	Action  string
	Service *Service
}

func Unmarshal(data []byte, v interface{}) error {
	switch value := v.(type) {
	case **registry.Service:
		log.Printf("**registry.Service")
		var service *Service
		err := json.Unmarshal(data, &service)
		if err != nil {
			return err
		}
		log.Printf("got %+v", service)
		log.Printf("got nodes %+v", *service.Nodes[0])
		*value = convertToV2Service(service)
		return nil

	case **registry.Result:
		log.Printf("**registry.Result")
		var result *Result
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		*value = convertToV2Result(result)
		log.Printf("got %+v", value)
		return nil
	}

	return fmt.Errorf("unknow type: %T", v)
}

func Marshal(v interface{}) ([]byte, error) {
	// data, err := json.Marshal(v)
	// return data, err
	switch value := v.(type) {
	case *registry.Service:
		data, err := json.Marshal(convertToV1Service(value))
		log.Printf("marshal %s", string(data))
		return data, err

	case *registry.Result:
		result := &Result{
			Action:  value.Action,
			Service: convertToV1Service(value.Service),
		}
		data, err := json.Marshal(result)
		log.Printf("marshal %s", string(data))
		return data, err
	}
	return nil, fmt.Errorf("unknow type: %T", v)
}

func convertToV2Service(service1 *Service) *registry.Service {
	// copy service
	service2 := &registry.Service{}
	service2.Name = service1.Name
	service2.Version = service1.Version
	service2.Metadata = service1.Metadata
	service2.Metadata = map[string]string{"protocol": "http"}
	service2.Endpoints = service1.Endpoints

	// copy nodes
	var nodes []*registry.Node
	for _, node := range service1.Nodes {
		address := node.Address
		if len(strings.Split(address, ":")) == 1 {
			address = fmt.Sprintf("%s:%d", address, node.Port)
		}
		nodes = append(nodes, &registry.Node{
			Id:       node.Id,
			Address:  address,
			Metadata: node.Metadata,
		})
	}
	service2.Nodes = nodes

	return service2
}

func convertToV2Result(result1 *Result) *registry.Result {
	return &registry.Result{
		Action:  result1.Action,
		Service: convertToV2Service(result1.Service),
	}
}

func convertToV1Service(s *registry.Service) *Service {
	result := &Service{
		Name:     s.Name,
		Version:  s.Version,
		Metadata: s.Metadata,
	}

	// copy nodes
	var nodes []*Node
	for _, node := range s.Nodes {
		address := node.Address
		port := 0
		strs := strings.Split(node.Address, ":")
		if len(strs) == 2 {
			address = strs[0]
			i, err := strconv.Atoi(strs[1])
			if err == nil {
				port = i
			}
		}
		if len(strs) == 1 {
			address = strs[0]
		}
		nodes = append(nodes, &Node{
			Id:       node.Id,
			Address:  address,
			Port:     port,
			Metadata: node.Metadata,
		})
	}
	result.Nodes = nodes

	// copy endpoints
	var eps []*registry.Endpoint
	for _, ep := range s.Endpoints {
		e := new(registry.Endpoint)
		*e = *ep
		eps = append(eps, e)
	}
	result.Endpoints = eps

	return result
}
