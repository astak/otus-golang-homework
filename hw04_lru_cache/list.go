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
	l.prepend(item)
	l.count++
	return
}

func (l *list) PushBack(v interface{}) (item *ListItem) {
	item = &ListItem{Value: v}
	l.append(item)
	l.count++
	return
}

func (l *list) Remove(i *ListItem) {
	l.detach(i)
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	l.detach(i)
	l.prepend(i)
}

func (l *list) prepend(i *ListItem) {
	if l.front == nil {
		l.front = i
		l.back = i
	} else {
		i.Next = l.front
		l.front.Prev = i
		l.front = i
	}
}

func (l *list) append(i *ListItem) {
	if l.back == nil {
		l.front = i
		l.back = i
	} else {
		i.Prev = l.back
		l.back.Next = i
		l.back = i
	}
}

func (l *list) detach(i *ListItem) {
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
}
