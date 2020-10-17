package json

import (
	"encoding/json"
	"fmt"
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
		var service *Service
		err := json.Unmarshal(data, &service)
		if err != nil {
			return err
		}
		*value = convertToV2Service(service)
		return nil

	case **registry.Result:
		var result *Result
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		*value = convertToV2Result(result)
		return nil
	}

	return fmt.Errorf("unknow type: %T", v)
}

func Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	return data, err
}

func convertToV2Service(service1 *Service) *registry.Service {
	// copy service
	result := &registry.Service{
		Name:      service1.Name,
		Version:   service1.Version,
		Metadata:  service1.Metadata,
		Endpoints: service1.Endpoints,
	}

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
	result.Nodes = nodes

	return result
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
