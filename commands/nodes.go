package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
)

func NodesExec(jenkins *gojenkins.Jenkins) error {
	nodes, err := jenkins.GetAllNodes()
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
	node.Poll()
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
