package queue

type node struct {
	Value int
	Next  *node
}
type Queue4Int struct {
	Head   *node
	Tail   *node
	Length int
}

func (q *Queue4Int) Push(value int) {
	n := node{
		Value: value,
	}
	if q.Length == 0 {
		q.Head = &n
		q.Tail = &n
	} else {
		q.Tail.Next = &n
		q.Tail = &n
	}
	q.Length++
}

func (q *Queue4Int) Pop() int {
	if q.Length == 0 {
		return -1
	}
	h := q.Head
	q.Head = h.Next
	q.Length--
	return h.Value
}

func NewQueue4Int() Queue4Int {
	q := Queue4Int{
		Length: 0,
	}
	return q
}
