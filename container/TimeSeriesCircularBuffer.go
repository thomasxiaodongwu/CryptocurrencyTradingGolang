package container

import "math"

const CAPACITY int = 128

type U struct {
	Time         int64
	Nearest_diff int64
	Item         interface{}
	Nearest      interface{}
}

// 这个可能是每秒是buffer_中的一个数值
type TimeSeriesCircularBuffer struct {
	capacity_ int
	size_     int
	index_    int
	full_     bool
	buffer_   []U
}

func (tscb *TimeSeriesCircularBuffer) clear() {
	tscb.size_ = 0
	tscb.index_ = 0
	tscb.full_ = false
	tscb.buffer_ = make([]U, tscb.capacity_+1)
}

func (tscb *TimeSeriesCircularBuffer) Empty() bool {
	return tscb.size_ == 0
}

func (tscb *TimeSeriesCircularBuffer) Full() bool {
	return tscb.full_
}

func (tscb *TimeSeriesCircularBuffer) Capacity() int {
	return tscb.capacity_ + 1
}

func (tscb *TimeSeriesCircularBuffer) Size() int {
	return tscb.size_
}

func (tscb *TimeSeriesCircularBuffer) Last_item(i int) U {
	return tscb.buffer_[(tscb.index_-i+tscb.capacity_)&tscb.capacity_]
}

func (tscb *TimeSeriesCircularBuffer) Last(i int) interface{} {
	return tscb.buffer_[(tscb.index_-i+tscb.capacity_)&tscb.capacity_].Item
}

func (tscb *TimeSeriesCircularBuffer) SetCapacity(capacity int) {
	if (capacity & (capacity - 1)) != 0 {
		panic("circular buffer capacity must be power of 2")
	}
	if capacity > CAPACITY {
		panic("circular buffer capacity must be less than 128")
	}
	tscb.capacity_ = capacity - 1
	tscb.clear()
}

func (tscb *TimeSeriesCircularBuffer) inc() {
	if !tscb.full_ {
		tscb.size_ += 1
	}
	// if (++index_ > capacity_) {
	if tscb.index_ > tscb.capacity_ {
		tscb.index_ = 0
		tscb.full_ = true
	}
	tscb.index_ += 1
}

// Time time, Time real, const T& item
// 时间向上取整
func (tscb *TimeSeriesCircularBuffer) PushBack(time int64, real int64, item interface{}) {
	if tscb.Empty() {
		tscb.buffer_[tscb.index_] = U{time, int64(math.Abs(float64(real - time))), item, item}
		tscb.inc()
	} else {
		if tscb.Last_item(0).Time == time {
			if int64(math.Abs(float64(real-time))) < tscb.Last_item(0).Nearest_diff {
				tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Nearest = item
			}
			tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Item = item
		} else {
			tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Item = tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Nearest

			var count int = int(time-tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Time)/1000 - 1
			for i := 0; i < count; i += 1 {
				tscb.buffer_[tscb.index_] = U{
					tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Time,
					0,
					tscb.buffer_[(tscb.index_+tscb.capacity_)&tscb.capacity_].Item,
					nil,
				}
				tscb.inc()
			}

			tscb.buffer_[tscb.index_] = U{
				time,
				int64(math.Abs(float64(real - time))),
				item,
				item,
			}
			tscb.inc()
		}
	}
}

// xli c++ version
//
// template<class T>
// class TimeSeriesCircularBuffer {
//     struct U {
//     Time time;
//     Time nearest_diff;
//     T item;
//     T nearest;
//     };
//     static const size_t CAPACITY = 128;
// public:
//     TimeSeriesCircularBuffer(size_t capacity = 128) {
//     set_capacity(capacity);
//     }

//     void set_capacity(size_t capacity) {
//     capacity_ = capacity - 1;
//     if ((capacity & (capacity - 1)) != 0) {
//         throw std::runtime_error("circular buffer capacity must be power of 2");
//     }
//     if (capacity > CAPACITY) {
//         throw std::runtime_error("circular buffer capacity must be less than " + std::to_string(CAPACITY));
//     }
//     }

//     // 时间向上取整
//     void push_back(Time time, Time real, const T& item) {
//     if (empty()) {
//         buffer_[index_] = U({time, abs(real-time), item, item});
//         inc();
//     } else {
//         U& l = const_cast<U&>(last_item());
//         if (l.time == time) {
//         if (abs(real-time) < l.nearest_diff) {
//             l.nearest = item;
//         }
//         l.item = item;
//         } else {
//         l.item = l.nearest;
//         auto count = (time - l.time) / NANO_PER_SEC - 1;
//         for(auto i = 0; i < count; ++i) {
//             buffer_[index_] = U({l.time, 0, l.item});
//             inc();
//         }
//         buffer_[index_] = U({time, abs(real-time), item, item});
//         inc();
//         }
//     }
//     }

//     void push_back(int time, T&& item) {
//     }

//     void push_back(const T& item) {
//     buffer_[index_] = {0, 0, item};
//     inc();
//     }

//     void push_back(T&& item) {
//     buffer_[index_] = {0, 0, std::move(item)};
//     inc();
//     }

//     const U& last_item() const {
//     return buffer_[(index_ + capacity_) & capacity_];
//     }

//     const T& last(size_t i = 0) const {
//     return buffer_[(index_ - i + capacity_) & capacity_].item;
//     }

//     size_t size() const {
//     return size_;
//     }

//     size_t capacity() const {
//     return capacity_ + 1;
//     }

//     bool empty() const {
//     return size_ == 0;
//     }

//     bool full() const {
//     return full_;
//     }

// private:
//     void inc() {
//     if (!full_) {
//         ++size_;
//     }

//     if (++index_ > capacity_) {
//         index_ = 0;
//         full_ = true;
//     }
//     }

//     size_t capacity_ = 0;
//     size_t size_ = 0;
//     size_t index_ = 0;
//     bool full_ = false;
//     template<typename R, size_t N>
//     struct __attribute__((aligned(64))) aligned_array : public std::array<U, N> {
//     };
//     aligned_array<U, CAPACITY> buffer_;
// };
