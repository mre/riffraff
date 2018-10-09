package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
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

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if online {
		fmt.Printf("%v %v: Online\n", green("✓"), node.GetName())
	} else {
		fmt.Printf("%v %v: Offline\n", red("✗"), node.GetName())
	}
	return nil
}
