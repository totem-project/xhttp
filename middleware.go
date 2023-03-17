package xhttp

import (
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Request Middleware(s)
//_______________________________________________________________________

func verifyRequestMethod(req *Request, c *Client) error {
	// req.Method in AllowMethods
	if req.RawRequest == nil {
		return errors.New("req.RawRequest is nil")
	}
	currentMethod := req.RawRequest.Method
	if funk.Contains(c.ClientOptions.AllowMethods, currentMethod) == false {
		return fmt.Errorf(`http method %s not allowed`, currentMethod)
	}
	return nil
}

func readRequestBody(req *Request, c *Client) error {
	_, err := req.GetBody()
	if err != nil {
		return err
	}
	return nil
}

func setTrace(req *Request) {
	// setTrace
	if req.clientTrace == nil && req.trace {
		req.clientTrace = &clientTrace{}
	}
	if req.clientTrace != nil {
		req.ctx = req.clientTrace.createContext(req.GetContext())
	}
}

func setRequestHeader(req *Request) {
	// fix header
	if req.RawRequest.Header.Get("Accept-Language") == "" {
		req.RawRequest.Header.Set("Accept-Language", "en")
	}
	if req.RawRequest.Header.Get("Accept") == "" {
		req.RawRequest.Header.Set("Accept", "*/*")
	}
}

func setContest(req *Request) {
	if req.GetContext() != nil {
		req.RawRequest = req.RawRequest.WithContext(req.GetContext())
	}
}

func createHTTPRequest(req *Request, c *Client) error {
	setTrace(req)
	setRequestHeader(req)
	setContest(req)

	// assign close connection option
	req.RawRequest.Close = c.ClientOptions.DisableKeepAlives
	// 更新 config 中的 headers cookie
	for key, value := range c.ClientOptions.Headers {
		// req 的 header 优先级高于默认设置的 header 值
		if len(req.RawRequest.Header.Values(key)) > 0 {
			continue
		}
		req.RawRequest.Header.Set(key, value)
	}
	if c.ClientOptions.Cookies != nil {
		for k, v := range c.ClientOptions.Cookies {
			req.RawRequest.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}
	return nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Response Middleware(s)
//_______________________________________________________________________

func readResponseBody(resp *Response, c *Client) error {
	// 限制响应长度
	lr := io.LimitReader(resp.RawResponse.Body, c.ClientOptions.MaxRespBodySize)
	bodyBytes, err := ioutil.ReadAll(lr)
	if err != nil {
		//if bodyBytes != nil {
		//	golog.Debugf("error when read http response body, but get bodyBytes: %v", err)
		//}
		return err
	}
	resp.Body = bodyBytes
	defer resp.RawResponse.Body.Close()
	return nil
}

func responseLogger(resp *Response, c *Client) error {
	if c.ClientOptions.Debug {
		req := resp.Request
		reqString, err := req.GetRaw()
		if err != nil {
			return err
		}

		respString, err := resp.GetRaw()
		if err != nil {
			return err
		}

		latency := resp.GetLatency()

		reqLog := "\n==============================================================================\n" +
			"--- REQUEST ---\n" +
			fmt.Sprintf("%s  %s  %s\n", req.GetMethod(), req.GetUrl().String(), req.RawRequest.Proto) +
			fmt.Sprintf("HOST   : %s\n", req.RawRequest.URL.Host) +
			fmt.Sprintf("RequestString:\n%s\n", reqString) +
			"------------------------------------------------------------------------------\n" +
			"--- RESPONSE ---\n" +
			fmt.Sprintf("STATUS       : %s\n", resp.RawResponse.Status) +
			fmt.Sprintf("PROTO        : %s\n", resp.RawResponse.Proto) +
			fmt.Sprintf("RECEIVED AT  : %v\n", resp.getReceivedAt().Format(time.RFC3339Nano)) +
			fmt.Sprintf("Attempt Num  : %d\n", req.attempt) +
			fmt.Sprintf("TIME DURATION: %v\n", latency) +
			fmt.Sprintf("HOST   : %s\n", req.RawRequest.URL.Host) +
			fmt.Sprintf("ResponseString:\n%s\n", respString) +
			"------------------------------------------------------------------------------\n"
		fmt.Println(reqLog)
	}
	return nil
}
