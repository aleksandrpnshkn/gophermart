package handlers

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestPing(t *testing.T) {
	apitest.New().
		HandlerFunc(Ping()).
		Get("/api/ping").
		Expect(t).
		Body(`{"result": true}`).
		Status(http.StatusOK).
		End()
}
