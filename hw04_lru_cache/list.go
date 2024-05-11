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
	front *ListItem
	back  *ListItem
	count int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) (item *ListItem) {
	item = &ListItem{Value: v}
	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}
	l.count++
	return
}

func (l *list) PushBack(v interface{}) (item *ListItem) {
	item = &ListItem{Value: v}
	if l.back == nil {
		l.back = item
		l.front = item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}
	l.count++
	return
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	i.Next = nil
	i.Prev = nil
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i == l.back {
		l.back = i.Prev
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	i.Next = l.front
	i.Prev = nil
	if l.front != nil {
		l.front.Prev = i
	}
	l.front = i
}
