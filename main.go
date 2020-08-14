package main

import (
	"fmt"
	"github.com/annchain/commongo/pq"
)

func main() {
	q := &pq.PriorityQueue{}
	q.Push(&pq.Item{
		Value:    5,
		Priority: 5,
	})
	q.Push(&pq.Item{
		Value:    3,
		Priority: 3,
	})
	q.Push(&pq.Item{
		Value:    7,
		Priority: 7,
	})
	fmt.Println(q.Len())
	fmt.Println(q.Pop().(*pq.Item).Value)
	fmt.Println(q.Pop().(*pq.Item).Value)
	fmt.Println(q.Pop().(*pq.Item).Value)
	fmt.Println(q.Pop().(*pq.Item).Value)
}
