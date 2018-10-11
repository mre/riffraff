package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
)

type Queue struct {
	jenkins *gojenkins.Jenkins
	regex   string
	verbose bool
	salt    bool
}

func NewQueue(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) *Queue {
	return &Queue{
		jenkins,
		regex,
		verbose,
		salt,
	}
}

func (q Queue) Exec() error {
	queue, err := q.jenkins.GetQueue()
	if err != nil {
		return err
	}
	fmt.Println(queue.Raw)
	return nil
}
