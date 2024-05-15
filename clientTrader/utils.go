package clientTrader

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"winter/global"
)

// Get 需求
func Get_Servers_Time() {
	proxy, _ := url.Parse(global.PROXY_URI)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Get("https://www.okex.com/api/v5/public/time")
	if err != nil {
		fmt.Println("出错了", err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
