package main

import (
	"net/http"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/assert"
)

func TestPing(t *testing.T) {
	handler := http.HandlerFunc(ping)
	ts := newTestServer(t, handler)
	defer ts.Close()

	res := ts.get(t, "/")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}
