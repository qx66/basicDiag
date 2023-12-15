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
	
	cli := nhttp.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := cli.Get(url)
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
	Id      string `json:"id,omitempty"`
}

func Report(url string, body []byte) (string, error) {
	var reportResponse ReportResponse
	
	cli := nhttp.Client{
		Timeout: 10 * time.Second,
	}
	
	// req
	resp, err := cli.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return reportResponse.Id, err
	}
	
	// req response body
	respBody := resp.Body
	defer respBody.Close()
	
	respBodyByte, err := io.ReadAll(respBody)
	if err != nil {
		return reportResponse.Id, err
	}
	
	if resp.StatusCode != 200 {
		return reportResponse.Id, errors.New(fmt.Sprintf("http code: %d, resp: %s.", resp.StatusCode, string(respBodyByte)))
	}
	
	err = json.Unmarshal(respBodyByte, &reportResponse)
	return reportResponse.Id, err
}
