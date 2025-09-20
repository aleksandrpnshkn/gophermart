package handlers

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestNotFound(t *testing.T) {
	apitest.New().
		HandlerFunc(NotFound()).
		Get("/404").
		Expect(t).
		Body(`{"error": {"message":"not found"}}`).
		Status(http.StatusNotFound).
		End()
}
