/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:48
 * @Last Modified by: xwu
 * @Last Modified time: 2022-06-01 13:16:36
 */
package operator

import (
	"math"
	"winter/container"
	"winter/mathpro"
)

type PointSkew struct {
	Data     *container.RingBufferFloat64
	Interval int64

	sum_    float64
	sum2_   float64
	sum3_   float64
	result_ float64
}

func (op *PointSkew) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval

	op.sum_ = 0
	op.sum2_ = 0
	op.sum3_ = 0
	op.result_ = 0
}

// x, x2
func (op *PointSkew) Update(time int64, x float64) float64 {
	// 剔除头部应该去掉的
	if op.Data.Full() {
		old_v := op.Data.Front()
		op.sum_ -= old_v
		op.sum2_ -= old_v * old_v
		op.sum3_ -= old_v * old_v * old_v
	}

	// 增加新的元素
	op.Data.Push_back(x)
	op.sum_ += x
	op.sum2_ += x * x
	op.sum3_ += x * x * x

	n := float64(op.Data.Len())

	if n < 3 {
		op.result_ = 0
	} else {
		var A float64 = op.sum_ / n
		var B float64 = op.sum2_/n - A*A
		var C float64 = op.sum3_/n - A*A*A - 3*A*B
		if B < 0.1 {
			op.result_ = 0
		} else {
			var R float64 = math.Sqrt(B)
			op.result_ = (math.Sqrt(n*(n-1)) * C) / ((n - 2) * R * B)
		}
	}

	if !mathpro.Isfinite(op.result_) {
		op.result_ = 0
	}
	return op.result_
}

func (op *PointSkew) Value() float64 {
	return op.result_
}

type TimeSkew struct {
	Data     []valueTm
	Interval int64

	sum_    float64
	sum2_   float64
	sum3_   float64
	result_ float64
}

func (op *TimeSkew) Init(interval int64) {
	op.Interval = interval

	op.sum_ = 0
	op.sum2_ = 0
	op.sum3_ = 0
	op.result_ = 0
}

// x, x2
func (op *TimeSkew) Update(time int64, x float64) float64 {
	new_value := valueTm{Tm: time, Value: x}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		old_v := op.Data[i].Value
		op.sum_ -= old_v
		op.sum2_ -= old_v * old_v
		op.sum3_ -= old_v * old_v * old_v
		i += 1
	}

	// 增加新的元素
	op.Data = append(op.Data[i:], new_value)
	op.sum_ += x
	op.sum2_ += x * x
	op.sum3_ += x * x * x

	n := float64(len(op.Data))

	if n < 3 {
		op.result_ = 0
	} else {
		var A float64 = op.sum_ / n
		var B float64 = op.sum2_/n - A*A
		var C float64 = op.sum3_/n - A*A*A - 3*A*B
		if B < 0.1 {
			op.result_ = 0
		} else {
			var R float64 = math.Sqrt(B)
			op.result_ = (math.Sqrt(n*(n-1)) * C) / ((n - 2) * R * B)
		}
	}

	if !mathpro.Isfinite(op.result_) {
		op.result_ = 0
	}
	return op.result_
}

func (op *TimeSkew) Value() float64 {
	return op.result_
}

// class OpSkew {
//     std::deque<double> data;
//     int intervals = 0;
//     double x1 = 0, x2 = 0, x3 = 0, value = 0;
//   public:
//     OpSkew(int intervals=5) : intervals(intervals) {
//     }

//     double update(double x) {
//       if (data.size() == intervals) {
//         remove(data.front());
//         data.pop_front();
//       }
//       data.push_back(x);
//       add(x);
//       compute();

//       return value;
//     }

//     double result() {
//       return value;
//     }

//   private:
//     void add(double x) {
//       auto v2 = x * x;
//       x1 += x;
//       x2 += v2;
//       x3 += v2 * x;
//     }

//     void remove(double x){
//       auto v2 = x * x;
//       x1 -= x;
//       x2 -= v2;
//       x3 -= v2 * x;
//     }

//     void compute() {
//       auto N = data.size();
//       if (N < 3) {
//         value = 0;
//       } else {
//         auto A = x1 / N;
//         auto B = x2 / N - A * A;
//         auto C = x3 / N - A * A * A - 3 * A * B;

//         if (almost_less_equal(B, 0.1)) {
//           value = 0;
//         } else {
//           auto R = sqrtf(B);
//           value = (sqrtf(N * (N - 1)) * C) / ((N - 2) * R * R * R);
//         }
//       }
//     }
//   };
