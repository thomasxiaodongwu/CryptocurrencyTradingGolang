/*
 * @Author: xwu
 * @Date: 2021-12-26 18:42:21
 * @Last Modified by:   xwu
 * @Last Modified time: 2021-12-26 18:42:21
 */
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"winter/global"
)

func OkexSign(timestamp int64) string {
	message := Strcat(
		strconv.FormatInt(timestamp, 10),
		"GET",
		"/users/self/verify",
	)

	h := hmac.New(sha256.New, []byte(global.Okex_SecertKey))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
