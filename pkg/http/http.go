package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	nhttp "net/http"
	"time"
)

// Get

type Result struct {
	StatusCode int
	Header     map[string][]string
	Body       string
	StartTime  int64
	EndTime    int64
}

func Get(url string) (Result, error) {
	var r Result
	r.StartTime = time.Now().Unix()
	
	resp, err := nhttp.Get(url)
	if err != nil {
		r.EndTime = time.Now().Unix()
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
	
	r.EndTime = time.Now().Unix()
	return r, nil
}

// report

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
