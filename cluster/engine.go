package cluster

import (
	"sync"
)

// Engine state define
type engineState int

const (
	//pending: engine added to cluster, but not been validated.
	statePending engineState = iota
	//unhealthy: engine is unreachable.
	stateUnhealthy
	//healthy: engine is ready reachable.
	stateHealthy
	//disconnected: engine is removed from discovery
	stateDisconnected
)

// Engine state mapping
var stateText = map[engineState]string{
	statePending:      "Pending",
	stateUnhealthy:    "Unhealthy",
	stateHealthy:      "Healthy",
	stateDisconnected: "Disconnected",
}

// Cluster is exported
type Engine struct {
	sync.RWMutex
	ID      string
	IP      string
	Name    string
	Addr    string
	Cpus    int64
	Memory  int64
	Version string

	//labels {"node","wh7", "storagedriver":"aufs", "kernelversion":"4.4.0" "operatingsystem":"centos6.8"}
	Labels map[string]string
	//images []*Image
	//containers map[string]*Container
	//client http.client //access to humpback-agent api, for engine-api call
	stopCh chan struct{}
	state  engineState
}

//NewEngine is exported
func NewEngine(ip string) *Engine {

	e := &Engine{
		IP:     ip,
		Labels: make(map[string]string),
		stopCh: make(chan struct{}),
		state:  stateUnhealthy,
	}
	return e
}

func (e *Engine) IsHealthy() bool {

	e.Lock()
	defer e.Unlock()
	return e.state == stateHealthy
}

func (e *Engine) setState(state engineState) {

	e.Lock()
	e.state = state
	e.Unlock()
}

func (e *Engine) Status() string {

	e.Lock()
	defer e.Unlock()
	return stateText[e.state]
}