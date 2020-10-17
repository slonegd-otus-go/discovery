package json

import (
	"encoding/json"
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

func Unmarshal(data []byte, v *registry.Service) error {
	v1 := convert(v)
	return json.Unmarshal(data, &v1)
}

func Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	return data, err
}

func convert(s *registry.Service) *Service {
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
