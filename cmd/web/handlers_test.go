package main

import (
	"net/http"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/assert"
)

func TestGetRoot(t *testing.T) {
	handler := http.HandlerFunc(getRoot)
	ts := newTestServer(t, handler)
	defer ts.Close()

	res := ts.get(t, "/")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "This is my website!")

}
