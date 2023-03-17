package xhttp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	diytesthttp "totem/v1/testutils/http"
)

func TestNewRequest(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := diytesthttp.CreateGenServer(t)
	defer ts.Close()

	restUrl := ts.URL + "/json-no-set"

	testMethod := MethodPost
	var testBody = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	xReq, _ := NewRequest(testMethod, restUrl, bytes.NewReader(testBody))

	require.Equal(t, restUrl, xReq.GetUrl().String(), "req.GetUrl id wrong")
	require.Equal(t, testMethod, xReq.GetMethod(), "req.GetMethod id wrong")
	requireBody, err := xReq.GetBody()
	require.Nil(t, err, "cannot use req.GetBody")
	require.Equal(t, testBody, requireBody, "req.GetBody id wrong")

	raw, err := xReq.GetRaw()
	require.Nil(t, err, "cannot use req.GetRaw")
	fmt.Println("1", string(raw))

	opt := DefaultClientOptions()
	opt.MaxRespBodySize = 100
	opt.Cookies = map[string]string{
		"key1":   "id1",
		"value1": "id2",
	}
	ctx := context.Background()
	resp, err := Do(ctx, xReq)
	require.Nil(t, err, "could not new http client")

	raw2, err := resp.Request.GetRaw()
	fmt.Println("2", string(raw2))
	//require.Contains(t, string(requireRaw), string(testBody), "raw is wrong")
}

func TestSetPostParam(t *testing.T) {
	xReq, _ := NewRequest("POST", "http://192.168.123.30:12345", nil)

	var params = map[string]string{
		"p1": "111",
		"p2": "222",
	}
	var paramList []string
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	p := strings.Join(paramList, "&")
	xReq.SetBody([]byte(p))
	raw, _ := xReq.GetRaw()
	fmt.Println(string(raw))
}
