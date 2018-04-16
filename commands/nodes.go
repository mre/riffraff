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
	waitGroup.Add(len(nodes))
	defer waitGroup.Wait()
	for _, node := range nodes {
		go printNodeStatus(&waitGroup, *node)
	}
	return nil
}

func printNodeStatus(waitGroup *sync.WaitGroup, node gojenkins.Node) error {
	defer waitGroup.Done()
	// Fetch Node Data
	_, err := node.Poll()
	if err != nil {
		return err
	}

	online, err := node.IsOnline()
	if err != nil {
		return err
	}
	if online {
		fmt.Printf("%v: Online\n", node.GetName())
	} else {
		fmt.Printf("%v: Offline\n", node.GetName())
	}
	return nil
}
