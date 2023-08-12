package main

import (
	"fmt"
	"time"
)

func main() {
	queue := TaskQueue{
		max_concurrent_tasks: 3,
		current_tasks:        0,
		index:                0,
		tasks: []Task{
			{
				id:       1,
				duration: 5,
			},
			{
				id:       2,
				duration: 2,
			},
			{
				id:       3,
				duration: 1,
			},
		},
		ch: make(chan string),
	}

	queue.process()

	queue.enqueue(Task{
		duration: 3,
	})

	queue.enqueue(Task{
		duration: 2,
	})

	for i := 0; i < len(queue.tasks); i++ {
		fmt.Println(<-queue.ch)
		queue.current_tasks--
		queue.process()
	}
}

type Task struct {
	id       int
	duration int
}

func (task Task) compute(ch chan string) {
	time.Sleep(time.Duration(task.duration) * time.Second)
	ch <- fmt.Sprintf("Task %d finished with result: %d", task.id, factorial(task.duration))
}

type TaskQueue struct {
	max_concurrent_tasks int
	current_tasks        int
	index                int
	tasks                []Task
	ch                   chan string
}

func (queue *TaskQueue) enqueue(task Task) {
	task.id = len(queue.tasks)
	queue.tasks = append(queue.tasks, task)
}

func (queue *TaskQueue) process() {
	for queue.index < len(queue.tasks) && queue.current_tasks < queue.max_concurrent_tasks {
		go queue.tasks[queue.index].compute(queue.ch)
		queue.current_tasks++
		queue.index++
	}
}

func factorial(iteration int) int {
	if iteration == 0 || iteration == 1 {
		return 1
	}

	return iteration * factorial(iteration-1)
}
