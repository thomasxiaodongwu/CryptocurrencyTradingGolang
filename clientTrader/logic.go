/*
 * @Author: xwu
 * @Date: 2022-05-23 16:20:25
 * @Last Modified by: xwu
 * @Last Modified time: 2022-07-15 14:01:34
 */
package clientTrader

import (
	"sync"
)

type LogicTrade struct {
	LDT LogicDoubleThreshold
}

type LogicDoubleThreshold struct {
	mutex       *sync.RWMutex
	mutexTarget *sync.RWMutex
	mutexParget *sync.RWMutex

	CurrentPosition int64
	TargetPosition  int64

	StartTm int64

	ThreEnter float64
	ThreExit  float64

	Cost float64

	Oid int64

	isTrading bool
}

func (LDT *LogicDoubleThreshold) Init(thre1 float64, thre2 float64) {
	LDT.isTrading = false
	LDT.CurrentPosition = 0
	LDT.TargetPosition = 0
	LDT.StartTm = 0
	LDT.ThreEnter = thre1
	LDT.ThreExit = thre2
	LDT.Oid = -1
	LDT.mutex = new(sync.RWMutex)
	LDT.mutexTarget = new(sync.RWMutex)
	LDT.mutexParget = new(sync.RWMutex)
}

// func (LDT *LogicDoubleThreshold) SetSpeed(temp string) {
// 	LDT.mutexSpeed.Lock() //写操作上锁
// 	LDT.tradeSpeed = temp
// 	LDT.mutexSpeed.Unlock()
// }

// func (LDT *LogicDoubleThreshold) GetSpeed() string {
// 	LDT.mutexSpeed.RLock() //写操作上锁
// 	tempSpeed := LDT.tradeSpeed
// 	LDT.mutexSpeed.RUnlock()
// 	return tempSpeed
// }

func (LDT *LogicDoubleThreshold) GetTradingStatus() bool {
	LDT.mutex.RLock()
	isTradingTemp := LDT.isTrading
	LDT.mutex.RUnlock()
	return isTradingTemp
}

func (LDT *LogicDoubleThreshold) TradingStatus(is_trading bool) {
	LDT.mutex.Lock() //写操作上锁
	LDT.isTrading = is_trading
	LDT.mutex.Unlock()
}

func (LDT *LogicDoubleThreshold) GetTargetPosition(tm int64, signal float64, ask1 float64, bid1 float64) int64 {
	var spread float64 = (ask1 - bid1) / (ask1 + bid1) * 20000
	// 第一版原来的逻辑是阈值开仓和衰减平仓
	target := LDT.GetCurrentPosition()
	if target == 0 {
		if signal > LDT.ThreEnter+spread {
			target = 1
		} else if signal < -1*LDT.ThreEnter-spread {
			target = -1
		}
	} else if target == -1 {
		if tm-LDT.StartTm > 3*60*1000 {
			target = 0
		} else {
			// adjust := math.Exp(-1 * float64(tm-LDT.StartTm-3000) / 60000)
			adjust := 0.
			if signal > LDT.ThreEnter+spread {
				target = 1
			} else if signal > -1*LDT.ThreExit*adjust+spread {
				target = 0
			}
		}
	} else if target == 1 {
		if tm-LDT.StartTm > 3*60*1000 {
			target = 0
		} else {
			// adjust := math.Exp(-1 * float64(tm-LDT.StartTm-3000) / 60000)
			adjust := 0.
			if signal < -LDT.ThreEnter-spread {
				target = -1
			} else if signal < LDT.ThreExit*adjust-spread {
				target = 0
			}
		}
	}

	LDT.setTargetPosition(target)
	return target
}

// func (LDT *LogicDoubleThreshold) GetTargetPosition(tm int64, signal float64, ask1 float64, bid1 float64) int64 {
// 	var spread float64 = (ask1 - bid1) / (ask1 + bid1) * 10000
// 	// 第一版原来的逻辑是阈值开仓和衰减平仓
// 	target := LDT.GetCurrentPosition()
// 	if target == 0 {
// 		if signal > LDT.ThreEnter+spread {
// 			target = 1
// 		}
// 	} else if target == 1 {
// 		adjust := 0.
// 		if signal < LDT.ThreExit*adjust-spread {
// 			target = 0
// 		}
// 	}

// 	LDT.setTargetPosition(target)
// 	return target
// }

// func (LDT *LogicDoubleThreshold) GetTargetPosition(tm int64, signal float64, ask1 float64, bid1 float64) int64 {
// 	var spread float64 = (ask1 - bid1) / (ask1 + bid1) * 20000
// 	// 第一版原来的逻辑是阈值开仓和衰减平仓
// 	target := LDT.GetCurrentPosition()
// 	if target == 0 {
// 		if signal < -1*LDT.ThreEnter-spread {
// 			target = -1
// 		}
// 	} else if target == -1 {
// 		adjust := 0.
// 		if signal > -1*LDT.ThreExit*adjust+spread {
// 			target = 0
// 		}
// 		// }
// 	}

// 	LDT.setTargetPosition(target)
// 	return target
// }

func (LDT *LogicDoubleThreshold) BindOid(oid int64) {
	// 这里有一个假设，那就是假设了从currentPosition变动到targetPosition,只需要一个订单就可以完成
	LDT.Oid = oid
}

func (LDT *LogicDoubleThreshold) UnBindOid(oid int64) bool {
	// 不管是cancel，还是成交，都需要把对应的仓位oid复位
	if oid == LDT.Oid {
		LDT.Oid = -1
		return true
	} else {
		return false
	}
}

func (LDT *LogicDoubleThreshold) GetCurrentPosition() int64 {
	LDT.mutexTarget.RLock()
	temp := LDT.CurrentPosition
	LDT.mutexTarget.RUnlock()
	return temp
}

func (LDT *LogicDoubleThreshold) setCurrentPosition(temp int64) {
	LDT.mutexTarget.Lock()
	LDT.CurrentPosition = temp
	LDT.mutexTarget.Unlock()
}

func (LDT *LogicDoubleThreshold) GetTargetPositionDirect() int64 {
	LDT.mutexParget.RLock()
	temp := LDT.TargetPosition
	LDT.mutexParget.RUnlock()
	return temp
}

func (LDT *LogicDoubleThreshold) setTargetPosition(temp int64) {
	LDT.mutexParget.Lock()
	LDT.TargetPosition = temp
	LDT.mutexParget.Unlock()
}

func (LDT *LogicDoubleThreshold) Active() bool {
	if LDT.TargetPosition != LDT.CurrentPosition {
		return true
	} else {
		return false
	}
}

func (LDT *LogicDoubleThreshold) IsBuy() bool {
	// TODO:这里其实有一个问题，那就是需要先判断Active，然后才能判断是买还是卖
	if LDT.TargetPosition > LDT.CurrentPosition {
		return true
	} else {
		return false
	}
}

func (LDT *LogicDoubleThreshold) IsOpen() (is_open bool, is_cross bool) {
	// TODO:这里其实有两个问题
	// 那就是需要先判断Active，然后才能判断是开仓还是平仓
	// 另一个是如果是穿越开仓的行为，那么就代表同时开仓和平仓了,然后如果穿越开仓，那么就都是先平再开
	// 判断是否是穿越开仓
	if LDT.CurrentPosition == -1 && LDT.TargetPosition == 1 {
		return false, true // 第一部分是平仓
	} else if LDT.CurrentPosition == 1 && LDT.TargetPosition == -1 {
		return false, true // 第一部分是平仓
	} else {
		if LDT.TargetPosition != 0 {
			return true, false // 只有一部分是开仓
		} else {
			return false, false // 只有一部分肯定是平仓
		}
	}
}

func (LDT *LogicDoubleThreshold) Update(avgprice float64, starttm int64) string {
	var ans string
	target := LDT.GetTargetPositionDirect()
	current := LDT.GetCurrentPosition()
	if target != 0 {
		LDT.StartTm = starttm
	} else {
		LDT.StartTm = 0
	}

	if target == 0 {
		LDT.Cost = 0
	} else {
		LDT.Cost = avgprice
	}

	if target > current {
		ans = "buy"
	} else if target < current {
		ans = "sell"
	}
	LDT.setCurrentPosition(target)
	return ans
}
