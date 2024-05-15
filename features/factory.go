/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:54
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-11 14:45:20
 */
package features

import "winter/params"

func NewAlphas(alphaName string, p params.AlphaParamsBase) Alphas {
	var ans Alphas
	// 这里尽量保持alphaName和结构体的名称是一致的
	switch alphaName {
	case "STA000":
		ans = new(STA000)
	default:
		panic("load signal " + alphaName + " failed.")
	}
	ans.Init(&p)
	return ans
}
