package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

//	type RequestParam struct {
//		Msg_type string            `json:"msg_type"`
//		Content  RequestParamChild `json:"content"`
//	}
//
//	type RequestParamChild struct {
//		Text string `json:"text"`
//	}
type HTTPRspBody struct {
	Result Results `json:"Result"`
}
type Results struct {
	RequestID     string   `json:"Result"`
	HasError      bool     `json:"HasError"`
	ResponseItems ErrorMsg `json:"ResponseItems"`
}
type ErrorMsg struct {
	ErrorMsg string `json:"ErrorMsg"`
}

func FeiShuMsg(message string) error {
	reqBody := `{"msg_type":"text","content":{"text":"` + message + `"}}`

	httpReq, err := http.NewRequest("POST", feishu_webhook, strings.NewReader(reqBody))
	if err != nil {
		logger.Info(Strcat(
			`NewRequest fail, feishu_webhook: `,
			feishu_webhook,
			`, reqBody: `,
			reqBody,
			`, err: `,
			err.Error(),
		))
		return err
	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		logger.Info(Strcat(
			`Do http fail, feishu_webhook: `,
			feishu_webhook,
			`, reqBody: `,
			reqBody,
			`, err: `,
			err.Error(),
		))
		return err
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	rspBody, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		logger.Info(Strcat(
			`ReadAll failed, feishu_webhook: `,
			feishu_webhook,
			`, reqBody: `,
			reqBody,
			`, err: `,
			err.Error(),
		))
		return err
	}

	if !strings.Contains(string(rspBody), "success") { // 没有日志就证明发送成功了
		logger.Info(`feishu message sented failed`)
	}

	// unmarshal: 解析HTTP返回的结果
	// body: {"Result":{"RequestId":"12131","HasError":true,"ResponseItems":{"ErrorMsg":"错误信息"}}}
	// {"StatusCode":0,"StatusMessage":"success"}

	// var result HTTPRspBody
	// if err := jsonIterator.Unmarshal(rspBody, &result); err != nil {
	// 	logger.Info(Strcat(
	// 		"Unmarshal fail, err: ",
	// 		err.Error(),
	// 	))
	// 	return err
	// }

	// if result.Result.HasError {
	// 	logger.Info(Strcat(
	// 		`http post fail, feishu_webhook: `,
	// 		feishu_webhook,
	// 		`, reqBody: `,
	// 		reqBody,
	// 		`, ErrorMsg:`,
	// 		result.Result.ResponseItems.ErrorMsg,
	// 	))
	// 	return errors.New(result.Result.ResponseItems.ErrorMsg)
	// }

	return nil

	// easy version
	// feishu_webhook := `https://open.feishu.cn/open-apis/bot/v2/hook/83da836d-8ff6-4de0-9806-8ca29db40300`

	// params := feishu_webhook.Values{}
	// params.Set("msg_type", "text")
	// params.Set("content", `{"text":"DataClient request example"}`)

	// logger.Info(params.Encode())
	// client := &http.Client{}

	// req, err := http.NewRequest("POST", feishu_webhook, strings.NewReader(`{"msg_type":"text","content":{"text":"DataClient request example"}}`))
	// if err != nil {
	// 	logger.Info("http.NewRequest failed")
	// }
	// // req.Header.Set("Content-Type", "application/json")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	logger.Info("error0")
	// }

	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Info("error1")
	// }

	// logger.Info(string(body))
}
