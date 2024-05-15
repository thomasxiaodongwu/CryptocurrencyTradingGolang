/*
 * @Author: xwu
 * @Date: 2021-12-26 18:42:27
 * @Last Modified by: xwu
 * @Last Modified time: 2022-01-22 16:40:51
 */
package utils

import "strings"

func Strcat(elems ...string) string {
	var n int = 0
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)

	for i := range elems {
		b.WriteString(elems[i])
	}
	return b.String()
}

func Strcatslice(elems []string) string {
	var n int = 0
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)

	for i := range elems {
		b.WriteString(elems[i])
	}
	return b.String()
}
