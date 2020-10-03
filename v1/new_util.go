package v1

import (
	"strconv"
	"strings"
	"fmt"

	"github.com/micro/go-micro/v2/registry"
)

func ConvertServiceToV2(service *Service) *registry.Service {
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

func ConvertServiceToV1(s *registry.Service) *Service {
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

func ConvertResultToV2(result *Result) *registry.Result {
	r := &registry.Result{
		Action  : result.Action,
		Service : ConvertServiceToV2(result.Service),
	}
	return r
}