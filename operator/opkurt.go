/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:10
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 14:40:04
 */
package operator

// class OpKurt {
//     std::deque<double> data;
//     int intervals = 0;
//     double x1 = 0, x2 = 0, x3 = 0, x4 = 0, value = 0;
//   public:
//     OpKurt(int intervals=5) : intervals(intervals) {
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
//       auto v3 = v2 * x;
//       auto v4 = v3 * x;
//       x1 += x;
//       x2 += v2;
//       x3 += v3;
//       x4 += v4;
//     }

//     void remove(double x) {
//       auto v2 = x * x;
//       auto v3 = v2 * x;
//       auto v4 = v3 * x;
//       x1 -= x;
//       x2 -= v2;
//       x3 -= v3;
//       x4 -= v4;
//     }

//     void compute() {
//       auto N = data.size();
//       if (N < 4) {
//         value = 0;
//       } else {
//         auto A = x1 / N;
//         auto R = A * A;

//         auto B = x2 / N - R;
//         R *= A;

//         auto C = x3 / N - R - 3 * A * B;
//         R *= A;

//         auto D = x4 / N - R - 6 * B * A * A - 4 * C * A;

//         if (almost_less_equal(B, 0.1)) {
//           value = 0;
//         } else {
//           auto K = (N * N - 1) * D / (B * B) - 3 * ((N - 1) * (N - 1));
//           value = K / ((N - 2) * (N - 3));
//         }
//       }
//     }
//   };
