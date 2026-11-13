package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type testServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *testServer {

	ts := httptest.NewServer(h)

	t.Cleanup(func() {
		ts.Close()
	})

	return &testServer{ts}
}

func (ts *testServer) fireRequest(t *testing.T, method string, path string, payload io.Reader, headers http.Header) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, payload)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range headers {
		req.Header[key] = value
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		res.Body.Close()
	})

	return res
}

func Document(t *testing.T, res *http.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func (ts *testServer) Get(t *testing.T, path string, headers http.Header) *http.Response {
	return ts.fireRequest(t, http.MethodGet, path, nil, headers)
}

func (ts *testServer) Post(t *testing.T, path string, payload io.Reader, headers http.Header) *http.Response {
	return ts.fireRequest(t, http.MethodPost, path, payload, headers)
}

func (ts *testServer) Put(t *testing.T, path string, payload io.Reader, headers http.Header) *http.Response {
	return ts.fireRequest(t, http.MethodPut, path, payload, headers)
}

func (ts *testServer) Patch(t *testing.T, path string, payload io.Reader, headers http.Header) *http.Response {
	return ts.fireRequest(t, http.MethodPatch, path, payload, headers)
}

func (ts *testServer) Delete(t *testing.T, path string, headers http.Header) *http.Response {
	return ts.fireRequest(t, http.MethodPatch, path, nil, headers)
}
