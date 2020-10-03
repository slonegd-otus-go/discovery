package nats

import (
	"encoding/json"
	"time"

	"github.com/micro/go-micro/v2/registry"
	"github.com/nats-io/nats.go"

	"github.com/slonegd-otus-go/discovery/registry/nats/convert"
)

type natsWatcher struct {
	sub *nats.Subscription
	wo  registry.WatchOptions
}

func (n *natsWatcher) Next() (*registry.Result, error) {
	var result *convert.Result // var result *registry.Result
	for {
		m, err := n.sub.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			continue
		} else if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(m.Data, &result); err != nil {
			return nil, err
		}
		if len(n.wo.Service) > 0 && result.Service.Name != n.wo.Service {
			continue
		}
		break
	}

	return convert.ResultToV2(result), nil // return result, nil
}

func (n *natsWatcher) Stop() {
	n.sub.Unsubscribe()
}
