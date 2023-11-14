package http

import (
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
