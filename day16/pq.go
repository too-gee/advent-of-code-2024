package main

import "github.com/too-gee/advent-of-code-2024/shared"

type State struct {
	loc  shared.Coord
	dir  string
	cost int
	path []shared.Coord
}

type PriorityQueue []*State

// Impliment a bunch of methods to satisfy the heap interface
func (q PriorityQueue) Len() int {
	return len(q)
}

func (q PriorityQueue) Less(i, j int) bool {
	return q[i].cost < q[j].cost
}

func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *PriorityQueue) Push(x interface{}) {
	item := x.(*State)
	*q = append(*q, item)
}

func (q *PriorityQueue) Pop() interface{} {
	last := len(*q) - 1
	item := (*q)[last]
	*q = (*q)[0:last]
	return item
}
