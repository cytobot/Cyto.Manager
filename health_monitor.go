package main

import (
	"time"

	pbd "github.com/cytobot/messaging/transport/discord"
	pb "github.com/cytobot/rpc/manager"
)

type healthMonitor struct {
	workerHealthStatus   map[string]*pbd.HealthCheckStatus
	listenerHealthStatus map[string]*pbd.HealthCheckStatus
}

func NewHealthMonitor() *healthMonitor {
	monitor := &healthMonitor{
		workerHealthStatus:   make(map[string]*pbd.HealthCheckStatus, 0),
		listenerHealthStatus: make(map[string]*pbd.HealthCheckStatus, 0),
	}

	monitor.startStatusCleanupInterval()

	return monitor
}

func (m *healthMonitor) UpdateWorkerHealthStatus(status *pbd.HealthCheckStatus) {
	m.workerHealthStatus[status.GetInstanceID()] = status
}

func (m *healthMonitor) UpdateListenerHealthStatus(status *pbd.HealthCheckStatus) {
	m.listenerHealthStatus[status.GetInstanceID()] = status
}

func (m *healthMonitor) GetAllWorkerStatus() []*pb.HealthCheckStatus {
	return getStatusList(m.workerHealthStatus)
}

func (m *healthMonitor) GetAllListenerStatus() []*pb.HealthCheckStatus {
	return getStatusList(m.listenerHealthStatus)
}

func getStatusList(statusMap map[string]*pbd.HealthCheckStatus) []*pb.HealthCheckStatus {
	list := make([]*pb.HealthCheckStatus, 0)
	for _, s := range statusMap {
		item := &pb.HealthCheckStatus{
			Timestamp:        s.Timestamp,
			InstanceID:       s.InstanceID,
			ShardID:          s.ShardID,
			Uptime:           s.Uptime,
			MemAllocated:     s.MemAllocated,
			MemSystem:        s.MemSystem,
			MemCumulative:    s.MemCumulative,
			TaskCount:        s.TaskCount,
			ConnectedServers: s.ConnectedServers,
			ConnectedUsers:   s.ConnectedUsers,
		}
		list = append(list, item)
	}
	return list
}

func (m *healthMonitor) startStatusCleanupInterval() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				cleanupOldStatus(m.workerHealthStatus)
				cleanupOldStatus(m.listenerHealthStatus)
			}
		}
	}()
}

func cleanupOldStatus(s map[string]*pbd.HealthCheckStatus) {
	now := time.Now().UTC()
	for key, ws := range s {
		statusTime := MapFromProtoTimestamp(ws.Timestamp)
		if now.Sub(statusTime) < (60 * time.Minute) {
			delete(s, key)
		}
	}
}
