package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"kuskcore/testutil"
)

func TestRPCCallJSON(t *testing.T) {
	requestBody := map[string]interface{}{
		"hello": "world",
	}

	// the TestResponse is same with api Response
	type TestResponse struct {
		Status string      `json:"status,omitempty"`
		Msg    string      `json:"msg,omitempty"`
		Data   interface{} `json:"data,omitempty"`
	}

	response := &TestResponse{}
	result := &TestResponse{
		Status: "success",
		Msg:    "right",
		Data:   requestBody,
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Inspect the request and ensure that it's what we expect.
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("got=%s; want=application/json", req.Header.Get("Content-Type"))
		}
		if !strings.HasPrefix(req.Header.Get("User-Agent"), "Kusk; ") {
			t.Errorf("got=%s; want prefix='Kusk; '", req.Header.Get("User-Agent"))
		}
		if req.URL.Path != "/example/rpc/path" {
			t.Errorf("got=%s want=/example/rpc/path", req.URL.Path)
		}
		un, pw, ok := req.BasicAuth()
		if !ok {
			t.Error("no user/password set")
		} else if un != "test-user" {
			t.Errorf("got=%s; want=test-user", un)
		} else if pw != "test-secret" {
			t.Errorf("got=%s; want=test-secret", pw)
		}

		decodedRequestBody := map[string]interface{}{}
		if err := json.NewDecoder(req.Body).Decode(&decodedRequestBody); err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()
		if !testutil.DeepEqual(decodedRequestBody, requestBody) {
			t.Errorf("got=%#v; want=%#v", decodedRequestBody, requestBody)
		}

		// Provide a dummy rpc response
		rw.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(result)
		if err != nil {
			t.Fatal(err)
		}
		rw.Write(data)
	}))
	defer server.Close()

	//response := map[string]string{}
	client := &Client{
		BaseURL:     server.URL,
		AccessToken: "test-user:test-secret",
	}
	err := client.Call(context.Background(), "/example/rpc/path", requestBody, &response)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the response is as we expect.
	if !testutil.DeepEqual(response, result) {
		t.Errorf(`got=%#v; want=%#v`, response, result)
	}

	// Ensure that supplying a nil response is OK.
	err = client.Call(context.Background(), "/example/rpc/path", requestBody, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRPCCallError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "a terrible error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL}
	wantErr := ErrStatusCode{URL: server.URL + "/error", StatusCode: 500}
	err := client.Call(context.Background(), "/error", nil, nil)
	if !testutil.DeepEqual(wantErr, err) {
		t.Errorf("got=%#v; want=%#v", err, wantErr)
	}
}

func TestCleanedURLString(t *testing.T) {
	u, _ := url.Parse("https://user:pass@foobar.com")
	want := "https://foobar.com"

	got := cleanedURLString(u)
	if got != want {
		t.Errorf("clean = %q want %q", got, want)
	}
}
