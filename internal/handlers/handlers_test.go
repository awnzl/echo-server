package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/awnzl/echo-server/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestNameHandler(t *testing.T) {
	want := nameResponse{
		ServiceName: "echo-server",
	}

	log  := logger.NewZap("info")
	defer log.Sync()
	handlers := New(log)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "localhost:8080", nil)

	handlers.nameHandler(w, req)

	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, err, err)

	got := nameResponse{}
	json.Unmarshal(data, &got)
	assert.Equal(t, want, got, "Incorrect response")
}

func TestEchoHandler(t *testing.T) {
	want := echoResponse{
		EchoWord: "test echo word",
	}

	log  := logger.NewZap("info")
	defer log.Sync()
	handlers := New(log)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST",
		"localhost:8080/echo",
		strings.NewReader(`{"word":"test echo word"}`),
	)

	handlers.echoHandler(w, req)

	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, err, err)

	got := echoResponse{}
	json.Unmarshal(data, &got)
	assert.Equal(t, want, got, "Incorrect response")
}
