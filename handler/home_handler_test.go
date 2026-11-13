package handler

import (
	"net/http"
	"testing"

	"github.com/antonisgkamitsios/verese/app"
	"github.com/go-playground/assert/v2"
)

func TestHome(t *testing.T) {
	application := app.NewTestApplication(t)
	ts := NewTestServer(t, New(application).SetupRoutes())

	res := ts.Get(t, "/", nil)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	doc := Document(t, res)
	title := doc.Find("h1").Text()
	assert.Equal(t, title, "Hello mafks")
}
