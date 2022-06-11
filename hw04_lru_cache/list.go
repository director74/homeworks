package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	cnt       int
	firstNode *ListItem
	lastNode  *ListItem
}

func (l *list) Len() int {
	return l.cnt
}

func (l *list) Front() *ListItem {
	return l.firstNode
}

func (l *list) Back() *ListItem {
	return l.lastNode
}

func (l *list) PushFront(v interface{}) *ListItem {
	oldFront := l.Front()
	item := ListItem{
		Value: v,
		Next:  oldFront,
		Prev:  nil,
	}

	l.firstNode = &item

	if oldFront == nil {
		l.lastNode = &item
	} else {
		oldFront.Prev = &item
	}

	l.cnt++

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Front() == nil {
		return l.PushFront(v)
	}

	oldBack := l.Back()
	item := ListItem{
		Value: v,
		Next:  nil,
		Prev:  oldBack,
	}

	oldBack.Next = &item
	l.lastNode = &item
	l.cnt++

	return &item
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev == nil {
		l.firstNode = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.lastNode = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.cnt--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
