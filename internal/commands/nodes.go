package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
)

type Nodes struct {
	jenkins *gojenkins.Jenkins
}

func NewNodes(jenkins *gojenkins.Jenkins) *Nodes {
	return &Nodes{
		jenkins,
	}
}

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
