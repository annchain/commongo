package todolist

import (
	"container/list"
	"github.com/sasha-s/go-deadlock"
	"time"
)

type Traceable interface {
	GetValue() interface{}
	GetId() string
}

type todoItem struct {
	value             Traceable
	tryTimes          int
	insertTime        time.Time
	lastExecutionTime time.Time
	executedEver      bool
	listItem          *list.Element
}

type TodoList struct {
	ExpireDuration          time.Duration
	MinimumIntervalDuration time.Duration
	MaxTryTimes             int

	items     *list.List
	todoItems map[string]*todoItem
	//mu        sync.RWMutex
	mu deadlock.RWMutex
}

func (t *TodoList) InitDefault() {
	t.items = list.New()
	t.todoItems = make(map[string]*todoItem)
}

func (t *TodoList) AddTask(value Traceable) {
	t.mu.Lock()
	defer t.mu.Unlock()

	itemId := value.GetId()
	if _, ok := t.todoItems[itemId]; ok {
		return
	}
	it := &todoItem{
		value:             value,
		tryTimes:          0,
		insertTime:        time.Now(),
		lastExecutionTime: time.Time{},
		executedEver:      false,
	}
	it.listItem = t.items.PushBack(it)
	t.todoItems[itemId] = it
}

func (t *TodoList) GetTask() (value Traceable) {
	// keep dequeue until matched task found
	t.mu.Lock()
	defer t.mu.Unlock()

	found := false
	for {
		// find one executable
		if len(t.todoItems) == 0 {
			return nil
		}
		// items should be sorted in lastExecutionTime desc
		firstNode := t.items.Front()
		firstItem := firstNode.Value.(*todoItem)

		if firstItem.executedEver && firstItem.insertTime.Add(t.ExpireDuration).Before(time.Now()) {
			// expire firstItem.
			t.removeTaskNoLock(firstItem.value)
			continue
		}
		if !firstItem.executedEver || firstItem.lastExecutionTime.Add(t.MinimumIntervalDuration).Before(time.Now()) {
			// use
			found = true
		}
		// whatever, break the loop so no blocking
		break
	}

	if !found {
		return nil
	}

	firstNode := t.items.Front()
	firstItem := firstNode.Value.(*todoItem)

	value = firstItem.value
	t.items.Remove(firstNode)

	// re-enqueue firstItem if tryTimes is ok
	if firstItem.tryTimes < t.MaxTryTimes {
		firstItem.executedEver = true
		firstItem.lastExecutionTime = time.Now()
		firstItem.tryTimes++
		firstItem.listItem = t.items.PushBack(firstItem)
	} else {
		// remove this task
		delete(t.todoItems, value.GetId())
	}

	return
}

func (t *TodoList) RemoveTask(value Traceable) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.removeTaskNoLock(value)
}

func (t *TodoList) removeTaskNoLock(value Traceable) {
	itemId := value.GetId()
	if v, ok := t.todoItems[itemId]; !ok {
		return
	} else {
		t.items.Remove(v.listItem)
		delete(t.todoItems, itemId)
	}
}

func (t *TodoList) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.todoItems)
}

func (t *TodoList) TaskIds() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var keys []string
	for k, _ := range t.todoItems {
		keys = append(keys, k)
	}
	return keys
}
