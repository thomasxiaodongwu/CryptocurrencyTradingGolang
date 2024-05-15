/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:33
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 23:27:01
 */
package operator

import (
	"winter/mathpro"
)

type PointDecay struct {
	m_intervals int64
	m_num       int64
	m_sum       float64
	m_diff      float64
	m_denom     float64 // must be 1
	m_queue     []float64
	m_idx       int64
	m_result    float64
}

func (op *PointDecay) Init(interval int64) {
	op.m_intervals = interval
	op.m_num = 0
	op.m_sum = 0
	op.m_diff = 0
	op.m_denom = 1
	op.m_queue = make([]float64, interval)
	op.m_idx = 0
	op.m_result = 0
}

// TODO: 这里的count不知道是不是会因为时间过长导致误差累计，所以可能需要每隔一段时间进行一次清空
func (op *PointDecay) Update(time int64, x float64) float64 {
	if !mathpro.Isfinite(x) {
		return op.m_result // Keep previous value
	}

	op.m_sum -= op.m_diff
	if op.m_num == op.m_intervals {
		op.m_num -= 1
		var tail int64 = op.m_idx + 1
		if tail > op.m_intervals-1 {
			tail = 0
		}
		op.m_diff -= op.m_queue[tail] / float64(op.m_intervals)
	}

	op.m_idx += 1
	if op.m_idx > op.m_intervals-1 {
		op.m_idx = 0
	}
	op.m_queue[op.m_idx] = x

	if mathpro.Isfinite(op.m_queue[op.m_idx]) {
		op.m_num += 1
		op.m_sum += op.m_queue[op.m_idx]
		op.m_diff += op.m_queue[op.m_idx] / float64(op.m_intervals)
	}

	op.m_result = op.m_sum / op.m_denom

	if op.m_num != op.m_intervals {
		op.m_denom += float64(op.m_intervals-op.m_num) / float64(op.m_intervals)
	}
	return op.m_result
}

func (op *PointDecay) Value() float64 {
	return op.m_result
}

type TimeDecay struct {
	Data []valueTm

	Interval int64

	size             int64
	coef             float64
	decay_weight_sum float64
	equal_weight_sum float64

	Result_ float64
}

func (op *TimeDecay) Init(interval int64) {
	op.Interval = interval
	op.size = 0
	op.coef = 0
	op.decay_weight_sum = 0
	op.equal_weight_sum = 0
	op.Result_ = 0
}

// TODO: 这里的count不知道是不是会因为时间过长导致误差累计，所以可能需要每隔一段时间进行一次清空
func (op *TimeDecay) Update(time int64, x float64) float64 {
	new_value := valueTm{Tm: time, Value: x}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		op.size -= 1
		op.coef -= float64(i + 1)
		op.decay_weight_sum -= op.Data[i].Value * float64(i+1)
		op.equal_weight_sum -= op.Data[i].Value
		i += 1
	}
	op.Data = op.Data[i:]
	op.coef -= float64(i) * float64(op.size)
	op.decay_weight_sum -= op.equal_weight_sum * float64(i)

	// 增加新的元素
	op.Data = append(op.Data, new_value)
	op.size += 1
	op.coef += float64(op.size)
	op.decay_weight_sum += float64(op.size) * x
	op.equal_weight_sum += x

	op.Result_ = op.decay_weight_sum / op.coef

	return op.Result_
}

func (op *TimeDecay) Value() float64 {
	return op.Result_
}

// class OpDecay {
// 	public:
// 	  OpDecay(int intervals = 5)
// 		: m_intervals(intervals),
// 		  m_num(0),
// 		  m_sum(0),
// 		  m_diff(0),
// 		  m_denom(1), // must be 1
// 		  m_queue(intervals, NAN),
// 		  m_idx(0),
// 		  m_result(0) {
// 	  }

// 	  float update(float x) {
// 		if (!std::isfinite(x)) return m_result; // Keep previous value

// 		m_sum -= m_diff;
// 		if (m_num == m_intervals) {
// 		  m_num -= 1;
// 		  int tail = m_idx + 1;
// 		  if (tail > m_intervals - 1) tail = 0;
// 		  m_diff -= m_queue[tail] / float(m_intervals);
// 		}

// 		++m_idx;
// 		if (m_idx > m_intervals - 1) m_idx = 0;
// 		m_queue[m_idx] = x;

// 		if (std::isfinite(m_queue[m_idx])) {
// 		  ++m_num;
// 		  m_sum += m_queue[m_idx];
// 		  m_diff += m_queue[m_idx] / float(m_intervals);
// 		}

// 		m_result = m_sum / m_denom;

// 		if (m_num != m_intervals) {
// 		  m_denom += float(m_intervals - m_num) / float(m_intervals);
// 		}
// 		return m_result;
// 	  }

// 	  int count() {
// 		return m_num;
// 	  }

// 	  float result() {
// 		return m_result;
// 	  }

// 	private:
// 	  int m_intervals;
// 	  int m_num;
// 	  float m_sum;
// 	  float m_diff;
// 	  float m_denom;
// 	  std::vector<float> m_queue;
// 	  int m_idx;
// 	  float m_result;
// 	};
