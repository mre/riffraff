package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
)

// Queue is a type that can be used to get information about a specific build queue.
type Queue struct {
	jenkins *gojenkins.Jenkins
	regex   string
	verbose bool
	salt    bool
}

// NewQueue will intialize a new Queue instance with the provided parameters.
func NewQueue(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) *Queue {
	return &Queue{
		jenkins,
		regex,
		verbose,
		salt,
	}
}

// Exec will get details about the calling build queue.
func (q Queue) Exec() error {
	queue, err := q.jenkins.GetQueue()
	if err != nil {
		return err
	}
	fmt.Println(queue.Raw)
	return nil
}
