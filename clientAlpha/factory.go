/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:48
 * @Last Modified by: xwu
 * @Last Modified time: 2022-09-30 20:10:11
 */
package clientAlpha

import (
	"reflect"
	"strconv"
	"winter/features"
	"winter/global"
	"winter/lightgbm"
	"winter/params"
)

func NewAlphaFactory() AlphaFactory {
	ans := AlphaFactory{}

	// 这个依旧是使用
	ans.AlphasExIDSymbols = make([]AlphaSymbol, global.FMCount)
	ans.TickTime = make([]int64, global.FMCount)
	ans.II = make([]map[string]int, global.AggParameters.Data.Count)
	for i := 0; i < global.AggParameters.Data.Count; i++ {
		ans.II[i] = global.AlphaUid[i]
	}

	instr_alpha := reflect.ValueOf(global.AggParameters.Alpha)
	for i := 0; i < global.FMCount; i++ {
		// 初始化 交易所+品种+因子
		// instr_data := reflect.ValueOf(global.AggParameters.Data)
		// instr_data_p := instr_data.FieldByName(global.ExIDList[i]).Interface().(params.DataParams)
		var exid string
		for k1 := range global.AlphaUid {
			for k2 := range global.AlphaUid[k1] {
				if global.AlphaUid[k1][k2] == i {
					exid = global.ExIDList[k1]
				}
			}
		}

		instr_alpha_p := instr_alpha.FieldByName(exid).Interface().(params.AlphaParams)

		ans.AlphasExIDSymbols[i].Factors = make([]features.Alphas, len(instr_alpha_p.P))
		ans.AlphasExIDSymbols[i].FactorCols = make([]string, len(instr_alpha_p.P))
		ans.AlphasExIDSymbols[i].MinAsk1Price = 1e13
		ans.AlphasExIDSymbols[i].MaxBid1Price = -1
		for f_i, f_v := range instr_alpha_p.P {
			ans.AlphasExIDSymbols[i].Factors[f_i] = features.NewAlphas(f_v.Name, f_v)
			ans.AlphasExIDSymbols[i].FactorCols[f_i] = f_v.Name + "-" + strconv.FormatInt(f_v.Interval, 10)
		}
	}

	// 载入lightgbm模型，使用的是网上开源的leaves
	useTransformation := true
	model, err := lightgbm.LGEnsembleFromFile("lgb.txt", useTransformation)
	if err != nil {
		panic(err)
	}
	ans.Model = model

	return ans
}
