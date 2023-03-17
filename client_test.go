package xhttp

import (
	"context"
	"github.com/kataras/golog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	http2 "totem/v1/testutils/http"
)

func TestClient_Do(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ctx := context.Background()

	want := "success"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()
		w.WriteHeader(200)
		w.Write([]byte(want))
	}))

	hr, _ := http.NewRequest("GET", testServer.URL, nil)
	req := &Request{
		RawRequest: hr,
	}
	resp, err := Do(ctx, req)
	require.Nil(t, err, "could not do request with context")
	require.Equal(t, want, string(resp.Body), "could not get correct resp Body")
}

func TestClient_Do_Redirect(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := http2.CreateRedirectServer(t)
	defer ts.Close()

	options := GetHTTPOptions()
	options.Headers = map[string]string{
		"user-agent": "aaa",
	}
	ctx := context.Background()

	xReq, _ := NewRequest("GET", ts.URL+"/redirect-1", nil)
	resp, err := DoWithRedirect(ctx, xReq)
	require.Nil(t, err, "could not do request with redirect")
	body := resp.GetBody()
	require.Equal(t, "<a href=\"/redirect-2\">Temporary Redirect</a>.\n\n", string(body), "could not use redirect client")

	resp1, err := Do(ctx, xReq)
	require.Nil(t, err, "could not do request with no redirect")
	body1 := resp1.GetBody()
	require.Equal(t, "<a href=\"/redirect-2\">Temporary Redirect</a>.\n\n", string(body1), "could not use redirect client")

}

func TestClient_Do_Cookie(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := http2.CreateRedirectServer(t)
	defer ts.Close()

	options := DefaultClientOptions()
	options.Headers = map[string]string{
		"user-agent": "aaa",
	}
	options.Cookies = map[string]string{
		"clientcookieid1": "id1",
		"clientcookieid2": "id2",
	}

	ctx := context.Background()

	xReq, _ := NewRequest("GET", ts.URL+"/redirect-1", nil)
	resp, err := Do(ctx, xReq)
	require.Nil(t, err, "could not do request with no redirect")
	require.Equal(t, "aaa", resp.Request.RawRequest.Header.Get("user-agent"))
	require.Equal(t, "clientcookieid1=id1; clientcookieid2=id2", resp.Request.RawRequest.Header.Get("cookie"))
}

func TestHeader_And_ResponseBodyLimit(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := http2.CreateGetServer(t)
	defer ts.Close()
	options := DefaultClientOptions()
	options.MaxRespBodySize = 100
	options.Cookies = map[string]string{
		"key1":   "id1",
		"value1": "id2",
	}

	ctx := context.Background()
	xReq, _ := NewRequest("GET", ts.URL+"/", nil)
	xReq.SetHeader("user-agent", "aaa")
	xReq.EnableTrace()
	resp, err := Do(ctx, xReq)
	require.Nil(t, err, "could not do request with client")
	require.Equal(t, "aaa", resp.Request.GetHeaders().Get("user-agent"))
	require.Equal(t, "key1=id1; value1=id2", resp.Request.GetHeaders().Get("cookie"))
}

func TestAutoGzip(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := http2.CreateGenServer(t)
	defer ts.Close()

	ctx := context.Background()

	testcases := []struct{ url, want string }{
		{ts.URL + "/gzip-test", "This is Gzip response testing"},
		{ts.URL + "/gzip-test-gziped-empty-Body", ""},
		{ts.URL + "/gzip-test-no-gziped-Body", ""},
	}
	for _, tc := range testcases {
		hr, _ := http.NewRequest("GET", tc.url, nil)
		req := &Request{RawRequest: hr}
		resp, err := Do(ctx, req)
		require.Nil(t, err, "could not do request")
		body := resp.GetBody()
		require.Equal(t, tc.want, string(body), "could not auto gzip")
	}
}

func TestTransportCookie(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := http2.CreateGetServer(t)
	defer ts.Close()
	ctx := context.Background()

	defaultRedirectClient.ClientOptions.Debug = true

	xReq, _ := NewRequest("GET", ts.URL+"/transport-cookie", nil)
	xReq.SetHeader("user-agent", "aaa")
	xReq.EnableTrace()

	testNum := 5

	for i := 0; i <= testNum; i++ {
		_, err = DoWithRedirect(ctx, xReq)
		golog.Info(xReq.RawRequest.Cookies())
		require.Nil(t, err, "could not do request with client")
	}
	require.Equal(t, "success4", xReq.RawRequest.Cookies()[testNum-1].Value, "could not transport cookie to multi client")
}
