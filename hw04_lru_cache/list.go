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
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return &list{
		len: 0,
	}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	if l.front != nil {
		l.front.Prev = &item
	} else {
		l.back = &item
	}
	l.front = &item
	l.len++

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{
		Value: v,
		Prev:  l.back,
		Next:  nil,
	}
	if l.back != nil {
		l.back.Next = &item
	} else {
		l.front = &item
	}
	l.back = &item
	l.len++
	return &item
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev
	next := i.Next
	switch {
	case i == l.back:
		l.back = i.Prev
		l.back.Next = nil
		i.Next = nil
		i.Prev = nil
		l.len--
	case i == l.front:
		l.front = i.Next
		l.front.Prev = nil
		i.Next = nil
		i.Prev = nil
		l.len--
	case prev.Next == i && next.Prev == i:
		prev.Next = next
		next.Prev = prev
		i.Next = nil
		i.Prev = nil
		l.len--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	prev := i.Prev
	next := i.Next
	switch {
	case i == l.front:
		return
	case i == l.back:
		l.back = prev
		l.back.Next = nil
		i.Prev = nil
		i.Next = l.front
		l.front.Prev = i
		l.front = i
		return
	case prev.Next == i && next.Prev == i:
		prev.Next = next
		next.Prev = prev
		i.Prev = nil
		i.Next = l.front
		l.front.Prev = i
		l.front = i
		return
	}
}
