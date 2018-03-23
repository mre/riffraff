package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
)

func QueueExec(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) error {
	queue, err := jenkins.GetQueue()
	if err != nil {
		return err
	}
	fmt.Println(queue.Raw)
	// for _, task := range tasks {
	// 	fmt.Println(task.GetWhy())
	// }
	return nil
}
