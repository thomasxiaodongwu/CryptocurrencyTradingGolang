/*
 * @Author: xwu
 * @Date: 2021-12-26 18:46:37
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-09 18:42:35
 */
package container

// TODO:可能需要改掉空接口，否则，直接使用slice的性能更好
type RingBuffer struct {
	size_ int64
	head_ int64
	tail_ int64
	rbuf_ []interface{}
}

func (r *RingBuffer) Push_back(x interface{}) {
	r.rbuf_[r.head_%r.size_] = x
	r.head_ += 1
}

func (r *RingBuffer) Get(k int64) interface{} {
	return r.rbuf_[(r.head_+k)%r.size_]
}

func (r *RingBuffer) Back() interface{} {
	return r.Get(-1)
}

func (r *RingBuffer) Front() interface{} {
	if r.Full() {
		return r.Get(0)
	} else {
		return r.rbuf_[0]
	}
}

func (r *RingBuffer) Full() bool {
	return r.head_ >= r.size_
}

func (r *RingBuffer) Len() int64 {
	return r.size_
}

func (r *RingBuffer) Clear() {
	for i := range r.rbuf_ {
		r.rbuf_[i] = 0
	}
	r.head_ = 0
	r.tail_ = 0
}

// TODO:可能需要改掉空接口，否则，直接使用slice的性能更好
type RingBufferInt64 struct {
	size_ int64
	head_ int64
	tail_ int64
	rbuf_ []int64
}

func (r *RingBufferInt64) Push_back(x int64) {
	r.rbuf_[r.head_%r.size_] = x
	r.head_ += 1
}

func (r *RingBufferInt64) Get(k int64) int64 {
	return r.rbuf_[(r.head_+k)%r.size_]
}

func (r *RingBufferInt64) Back() int64 {
	return r.Get(-1)
}

func (r *RingBufferInt64) Front() int64 {
	if r.Full() {
		return r.Get(0)
	} else {
		return r.rbuf_[0]
	}
}

func (r *RingBufferInt64) Full() bool {
	return r.head_ >= r.size_
}

func (r *RingBufferInt64) Len() int64 {
	return r.size_
}

func (r *RingBufferInt64) Clear() {
	for i := range r.rbuf_ {
		r.rbuf_[i] = 0
	}
	r.head_ = 0
	r.tail_ = 0
}

// TODO:可能需要改掉空接口，否则，直接使用slice的性能更好
type RingBufferFloat64 struct {
	size_ int64
	head_ int64
	tail_ int64
	rbuf_ []float64
}

func (r *RingBufferFloat64) Push_back(x float64) {
	r.rbuf_[r.head_%r.size_] = x
	r.head_ += 1
}

func (r *RingBufferFloat64) Get(k int64) float64 {
	return r.rbuf_[(r.head_+k+r.size_)%r.size_]
}

func (r *RingBufferFloat64) Back() float64 {
	return r.Get(-1)
}

func (r *RingBufferFloat64) Front() float64 {
	if r.Full() {
		return r.Get(0)
	} else {
		return r.rbuf_[0]
	}
}

func (r *RingBufferFloat64) Full() bool {
	return r.head_ >= r.size_
}

func (r *RingBufferFloat64) Len() int64 {
	if r.head_ >= r.size_ {
		return r.size_
	} else {
		return r.head_
	}
}

func (r *RingBufferFloat64) Clear() {
	for i := range r.rbuf_ {
		r.rbuf_[i] = 0
	}
	r.head_ = 0
	r.tail_ = 0
}
