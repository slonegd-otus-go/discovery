package json

import (
	"encoding/json"
	"fmt"
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
	switch v.(type) {
	case *registry.Service:
		var service *Service
		err := json.Unmarshal(data, &service)
		if err != nil {
			return err
		}
		v = serviceToV2(service)
		return nil

	case *registry.Result:
		var result *Result
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		v = resultToV2(result)
		return nil
	}

	return fmt.Errorf("unknow type: %T", v)
}

func Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	return data, err
}

func serviceToV2(service *Service) *registry.Service {
	// copy service
	s := &registry.Service{
		Name:     service.Name,
		Version:  service.Version,
		Metadata: service.Metadata,
	}

	// copy nodes
	var nodes []*registry.Node
	for _, node := range service.Nodes {
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
	s.Nodes = nodes

	// copy endpoints
	var eps []*registry.Endpoint
	for _, ep := range service.Endpoints {
		e := new(registry.Endpoint)
		*e = *ep
		eps = append(eps, e)
	}
	s.Endpoints = eps

	return s
}

func resultToV2(result *Result) *registry.Result {
	r := &registry.Result{
		Action:  result.Action,
		Service: serviceToV2(result.Service),
	}
	return r
}
