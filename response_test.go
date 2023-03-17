package xhttp

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"totem/v1/testutils/http"
)

func TestNewResponse(t *testing.T) {
	err := Init()
	require.Nil(t, err, "could not do init xhttp pkg")
	require.NotEmpty(t, defaultClient, "could not make defaultClient")
	require.NotEmpty(t, defaultRedirectClient, "could not make defaultRedirectClient")

	ts := http.CreateGetServer(t)
	defer ts.Close()

	options := DefaultClientOptions()
	options.Cookies = map[string]string{
		"key1":   "id1",
		"value1": "id2",
	}
	ctx := context.Background()

	xReq, _ := NewRequest("GET", ts.URL+"/", nil)
	resp, err := Do(ctx, xReq)
	require.Nil(t, err)
	body := resp.GetBody()
	require.Nil(t, err)
	require.Equal(t, string(body), "TestGet: text response")
	flag := false
	if resp.GetLatency() > 5 {
		flag = true
	}
	require.Equal(t, flag, true)
}
