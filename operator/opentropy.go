/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:17
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 14:39:33
 */
package operator

// class OpEntropy {
// absl::btree_map<double, std::pair<int, double>> counts;
// std::deque<double> data;
// int intervals = 0;
// double value = 0;
// public:
// OpEntropy(int intervals=5) : intervals(intervals) {
// }

// double update(double x) {
// 	double old = 0;
// 	bool exist_old = false;
// 	if (data.size() == intervals) {
// 	exist_old = true;
// 	old = data.front();
// 	data.pop_front();
// 	}
// 	data.push_back(x);

// 	if (almost_equal(old, x)) {
// 	return value;
// 	} else {
// 	// old
// 	if(exist_old) {
// 		auto& pair = counts[old];
// 		auto count = pair.first;
// 		auto summand = pair.second;
// 		value += summand;

// 		if (count == 1) {
// 		counts.erase(old);
// 		} else {
// 		auto p_old = double(count - 1) / intervals;
// 		auto log_p_old = log2(p_old);
// 		pair.first = count - 1;
// 		pair.second = p_old * log_p_old;
// 		value -= pair.second;
// 		}
// 	}

// 	// new
// 	{
// 		auto& pair = counts[x];
// 		pair.first++;
// 		auto summand = pair.second;
// 		auto p_new = double(pair.first) / intervals;
// 		auto log_p_new = log2(p_new);
// 		pair.second = p_new * log_p_new;
// 		value += summand - pair.second;
// 	}
// 	}

// 	return value;
// }

// double result() {
// 	return value;
// }
// };
