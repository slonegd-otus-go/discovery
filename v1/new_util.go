package v1

import (
	"strconv"
	"strings"

	"github.com/micro/go-micro/v2/registry"
)

// func Copy(current []*Service) []*registry.Service {
// 	var services []*registry.Service

// 	for _, service := range current {
// 		// copy service
// 		s := &registry.Service{
// 			Name:     service.Name,
// 			Version:  service.Version,
// 			Metadata: service.Metadata,
// 		}

// 		// copy nodes
// 		var nodes []*registry.Node
// 		for _, node := range service.Nodes {
// 			nodes = append(nodes, &registry.Node{
// 				Id:       node.Id,
// 				Address:  fmt.Sprintf("%s:%d", node.Address, node.Port),
// 				Metadata: node.Metadata,
// 			})
// 		}
// 		s.Nodes = nodes

// 		// copy endpoints
// 		var eps []*registry.Endpoint
// 		for _, ep := range service.Endpoints {
// 			e := new(registry.Endpoint)
// 			*e = *ep
// 			eps = append(eps, e)
// 		}
// 		s.Endpoints = eps

// 		// append service
// 		services = append(services, s)
// 	}

// 	return services
// }

func ConvertToV1(s *registry.Service) *Service {
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
