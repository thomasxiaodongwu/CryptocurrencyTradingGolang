package container

type RingBufferGeneric[T any] struct {
	head_, tail_, size_, capacity_ int

	interval_ int64
	buffer_   []T
}

func (rb *RingBufferGeneric[T]) IsFull() bool {
	return rb.size_ == rb.capacity_
}

func (rb *RingBufferGeneric[T]) IsEmpty() bool {
	return rb.size_ == 0
}

func (rb *RingBufferGeneric[T]) Size() int {
	return rb.size_
}

func (rb *RingBufferGeneric[T]) Front() T {
	return rb.buffer_[rb.head_]
}

func (rb *RingBufferGeneric[T]) Back() T {
	return rb.buffer_[(rb.tail_-1)&(rb.capacity_-1)]
}

func (rb *RingBufferGeneric[T]) inc_() {
	rb.head_ = (rb.head_ + 1) & (rb.capacity_ - 1)
	rb.size_ -= 1
}

func (rb *RingBufferGeneric[T]) inc() {
	rb.tail_ = (rb.tail_ + 1) & (rb.capacity_ - 1)
	rb.size_ += 1
}

func (rb *RingBufferGeneric[T]) PushBack(tm int64, item T) {
	if rb.IsFull() {
		// double the capacity
		rb.buffer_ = append(rb.buffer_, rb.buffer_...)
		if rb.tail_ < rb.head_ {
			rb.head_ += rb.capacity_
		}
		rb.capacity_ = 2 * rb.capacity_
	}

	// fmt.Println(rb.head_, rb.tail_, rb.size_)
	rb.buffer_[rb.tail_] = item
	// fmt.Println(rb.time_)

	rb.inc()
}

func (rb *RingBufferGeneric[T]) PopFront() (t T) {
	if rb.size_ > 0 {
		ans := rb.buffer_[rb.head_]
		rb.inc_()
		return ans
	} else {
		logger.Info("pop from a empty ringbuffer")
		return t
	}
}

func (rb *RingBufferGeneric[T]) Clear() {
	rb.head_ = 0
	rb.tail_ = 0
	rb.size_ = 0
}
