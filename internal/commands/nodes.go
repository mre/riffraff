package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
)

// Nodes is a type used to request the status of different nodes used by a Jenkins server.
type Nodes struct {
	jenkins *gojenkins.Jenkins
}

// NewNodes will construct a new Nodes instance with the provided Jenkins client that
// has previously been configured with valid credentials.
func NewNodes(jenkins *gojenkins.Jenkins) *Nodes {
	return &Nodes{
		jenkins,
	}
}

// Exec will request all nodes connected to the configured Jenkins server and
// prints the status of each node.
func (n Nodes) Exec() error {
	nodes, err := n.jenkins.GetAllNodes()
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup

	for i := range nodes {
		waitGroup.Add(1)
		go func(node *gojenkins.Node) {
			defer waitGroup.Done()
			print(node)
		}(nodes[i])
	}

	waitGroup.Wait()
	return nil
}

func print(node *gojenkins.Node) error {

	// Fetch Node Data
	_, err := node.Poll()
	if err != nil {
		return err
	}

	online, err := node.IsOnline()
	if err != nil {
		return err
	}

	status := "Unknown"
	marker := Unknown
	if online {
		status = "Online"
		marker = Good
	} else {
		status = "Offline"
		marker = Bad
	}

	fmt.Printf("%s %s: %s\n", marker, node.GetName(), status)

	return nil
}
