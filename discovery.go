// копипаст из git.sedmax.ru/KIT/interservice - не работает с модулями
package discovery

import (
	"fmt"

	"github.com/micro/go-micro/v2/client/selector"
)

type Discovery struct {
	selector selector.Selector
}

func New(selector selector.Selector) *Discovery {
	return &Discovery{selector}
}

func (discovery *Discovery) GetAddress(serviceName string) (string, error) {
	next, err := discovery.selector.Select(serviceName)
	if err != nil {
		return "", fmt.Errorf("got function for next node of service %s failed: %s", serviceName, err)
	}
	node, err := next()
	if err != nil {
		return "", fmt.Errorf("got the next node for service %s failed: %s", serviceName, err)
	}

	addr := node.Address
	// if node.Port > 0 {
	// 	addr = fmt.Sprintf("http://%s:%d", addr, node.Port)
	// }
	return addr, nil
}
