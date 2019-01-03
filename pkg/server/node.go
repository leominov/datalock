package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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
	n.mu.Lock()
	state := n.IsHealthy()
	n.mu.Unlock()
	n.Healthy = state
	return n
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
	cli := http.DefaultClient
	cli.Timeout = 5 * time.Second
	resp, err := cli.Get(addr)
	if err != nil {
		log.Printf("Error requesting node %s state: %s", n.NodeName, err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Node %s has incorrect status: %s", n.NodeName, resp.Status)
		return false
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading node %s state: %s", n.NodeName, err)
		return false
	}
	state := strings.TrimSpace(string(b))
	if state == "ok" {
		return true
	}
	log.Printf("Node %s has incorrect state: %s", n.NodeName, state)
	return false
}
