/*
 * @Author: xwu
 * @Date: 2021-12-26 18:46:29
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-31 22:53:43
 */
package container

func NewRingBuffer(size int64) *RingBuffer {
	r := new(RingBuffer)
	r.size_ = size
	r.rbuf_ = make([]interface{}, size)
	r.head_ = 0
	r.tail_ = 0
	return r
}

func NewRingBufferInt64(size int64) *RingBufferInt64 {
	r := new(RingBufferInt64)
	r.size_ = size
	r.rbuf_ = make([]int64, size)
	r.head_ = 0
	r.tail_ = 0
	return r
}

func NewRingBufferFloat64(size int64) *RingBufferFloat64 {
	r := new(RingBufferFloat64)
	r.size_ = size
	r.rbuf_ = make([]float64, size)
	r.head_ = 0
	r.tail_ = 0
	return r
}

// // New creates a new skip list with default parameters. Returns a pointer to the new list.
// func NewSkipList(opt ...int64) *SkipList {
// 	var level = DEFAULT_SKIP_LIST_LEVEL
// 	if opt != nil {
// 		level = opt[0]
// 	}
// 	return NewWithMaxLevel(int(level))
// }

// TODO:还没有完成
func NewTimeSeriesCircularBuffer(capacity int) *TimeSeriesCircularBuffer {
	ans := new(TimeSeriesCircularBuffer)
	ans.SetCapacity(capacity)
	return ans
}
