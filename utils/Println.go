/*
 * @Author: xwu
 * @Date: 2022-01-22 21:52:48
 * @Last Modified by: xwu
 * @Last Modified time: 2022-01-22 21:57:09
 */
package utils

import (
	"fmt"
	"winter/params"
)

func PrintHftStrategyParam(ExID int, strategy_name string, x params.HftStrategyParam) {
	// for threshold_ii := range x.ThreEnter {
	// 	a := strconv.FormatFloat(
	// 		x.ThreEnter[threshold_ii],
	// 		'f',
	// 		1,
	// 		64,
	// 	)
	// 	b := strconv.FormatFloat(
	// 		x.ThreExit[threshold_ii],
	// 		'f',
	// 		1,
	// 		64,
	// 	)
	// 	logger.Info(Strcat(
	// 		"Exchange:",
	// 		global.ExIDList[ExID],
	// 		" ",
	// 		strategy_name,
	// 		" | Symbol:",
	// 		x.Symbol,
	// 		" | Threshold_Enter:",
	// 		a,
	// 		" | Threshold_Exit:",
	// 		b,
	// 	))
	// }
	fmt.Println("waiting for adapting for change")
}
