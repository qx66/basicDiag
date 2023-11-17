package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	nhttp "net/http"
)

type Result struct {
	StatusCode int
	Header     map[string][]string
	Body       string
}

func Get(url string) (Result, error) {
	var r Result
	
	resp, err := nhttp.Get(url)
	if err != nil {
		return r, err
	}
	
	r.StatusCode = resp.StatusCode
	r.Header = resp.Header
	
	respBody := resp.Body
	defer respBody.Close()
	
	// 请求成功之后，后续err都不影响返回的error
	respBodyByte, err := io.ReadAll(respBody)
	if err != nil {
		r.Body = fmt.Sprintf("read resp error, err: %s", err.Error())
	} else {
		r.Body = string(respBodyByte)
	}
	return r, nil
}

type ReportResponse struct {
	ErrCode int    `json:"errCode,omitempty"`
	ErrMsg  string `json:"errMsg,omitempty"`
	Result  string `json:"result,omitempty"`
}

func Report(url string, body []byte) (string, error) {
	var reportResponse ReportResponse
	
	// req
	resp, err := nhttp.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return reportResponse.Result, err
	}
	
	if resp.StatusCode != 200 {
		return reportResponse.Result, errors.New(fmt.Sprintf("http code: %d", resp.StatusCode))
	}
	
	// req response body
	respBody := resp.Body
	defer respBody.Close()
	
	respBodyByte, err := io.ReadAll(respBody)
	if err != nil {
		return reportResponse.Result, err
	}
	
	err = json.Unmarshal(respBodyByte, &reportResponse)
	return reportResponse.Result, err
}
