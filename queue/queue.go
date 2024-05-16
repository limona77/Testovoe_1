package queue

import (
	"Testovoe/client"
	"fmt"
)

type Queue struct {
	queue []client.Client
}

func (q *Queue) Enqueue(element client.Client) []client.Client {
	q.queue = append(q.queue, element)

	return q.queue
}

func (q *Queue) Dequeue() client.Client {
	element := q.queue[0]
	q.queue = q.queue[1:]
	return element
}

func (q *Queue) Front() (client.Client, error) {
	if len(q.queue) == 0 {
		return client.Client{}, fmt.Errorf("queue is empty")
	}

	return q.queue[0], nil
}

func (q *Queue) IsEmpty() bool {
	if len(q.queue) == 0 {
		return true
	} else {
		return false
	}
}

func (q *Queue) Print() {
	for _, item := range q.queue {
		fmt.Print(item, " ")
	}
	fmt.Println()
}

func (q *Queue) Len() int {
	return len(q.queue)
}
