package main

import (
	"fmt"
	"github.com/annchain/commongo/pq"
	"github.com/annchain/commongo/todolist"
	"time"
)

func pqTest() {
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

type Task struct {
	value int
}

func (t Task) GetValue() interface{} {
	return t.value
}

func (t Task) GetId() string {
	return fmt.Sprintf("ID-%d", t.value)
}

func todolistTest() {
	l := todolist.TodoList{
		ExpireDuration:          time.Second * 10,
		MinimumIntervalDuration: time.Second * 5,
		MaxTryTimes:             3,
	}
	l.InitDefault()
	for i := 0; i < 10; i++ {
		l.AddTask(Task{
			value: i,
		})
	}

	for {
		v := l.GetTask()
		fmt.Println("deque", v)
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	todolistTest()
}
