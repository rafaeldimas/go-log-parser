package queue

type Queue[T any] interface {
	Enqueue(T)
	Dequeue() T
	IsEmpty() bool
	Length() int
}

type queue[T any] struct {
	items []T
}

func New[T any]() Queue[T] {
	return &queue[T]{}
}

func (q *queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *queue[T]) Dequeue() T {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *queue[T]) Length() int {
	return len(q.items)
}
