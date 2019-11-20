package main

import (
	"encoding/json"
)

func (manager *managerState) setupHealthCheckListener() {
	// TODO: Keep track of individual shard state for cluster reporting.

	listenerHealth := manager.nats.Subscribe("listener_health")
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-listenerHealth:
				healthMsg := &HealthMessage{}
				json.Unmarshal(msg.Data, healthMsg)

				//log.Printf("[Discord Health] %+v", healthMsg)
			case <-quit:
				return
			}
		}
	}()
}
