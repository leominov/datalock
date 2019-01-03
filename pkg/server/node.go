package server

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/leominov/datalock/pkg/util/httpget"
)

type Node struct {
	Root     bool
	NodeName string
	Hostname string
	mu       sync.Mutex
	Healthy  bool
}

func NodeFromAddr(addr string, root bool) *Node {
	n := &Node{
		Root:     root,
		NodeName: addr,
		Hostname: addr,
	}
	if n.Root {
		n.NodeName = "root"
	}
	state := n.IsHealthy()
	n.SetHealthy(state)
	return n
}

func (n *Node) SetHealthy(state bool) {
	n.mu.Lock()
	n.Healthy = state
	n.mu.Unlock()
}

func (n *Node) State() string {
	if n.Healthy {
		return "healthy"
	}
	return "inactive"
}

func (n *Node) String() string {
	if n.NodeName == n.Hostname {
		return n.NodeName
	}
	return fmt.Sprintf("%s (%s)", n.NodeName, n.Hostname)
}

func (n *Node) AbsoluteLink(link string) string {
	return fmt.Sprintf(SeriesLinkFormat, n.Hostname, link)
}

func (n *Node) IsHealthy() bool {
	if n.Root {
		// Root node is always healthy
		return true
	}
	addr := fmt.Sprintf("http://%s/healthz", n.Hostname)
	body, err := httpget.HttpGet(addr)
	if err != nil {
		log.Printf("Error getting node %s state: %s", n.NodeName, err)
		return false
	}
	state := strings.TrimSpace(body)
	if state == "ok" {
		return true
	}
	log.Printf("Node %s has incorrect state: %s", n.NodeName, state)
	return false
}
