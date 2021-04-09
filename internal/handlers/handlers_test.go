package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/awnzl/echo-server/internal/logger"
)

type testSuite struct{
	suite.Suite
	log      *zap.Logger
	handlers *Handlers
}

func (s *testSuite) SetupTest() {
	s.log = logger.NewZap("info")
	s.handlers = New(s.log)
}

func (s *testSuite) TearDownSuite() {
	s.log.Sync()
}

func (s *testSuite) TestNameHandler() {
	want := `{"service_name":"echo-server"}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "localhost:8080", nil)
	s.handlers.nameHandler(w, req)

	res := w.Result()
	got, err := io.ReadAll(res.Body)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.JSONEq(want, string(got), "Incorrect response")
}

func (s *testSuite) TestEchoHandler() {
	want := `{"echo":"test echo word"}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST",
		"localhost:8080/echo",
		strings.NewReader(`{"word":"test echo word"}`),
	)
	s.handlers.echoHandler(w, req)

	res := w.Result()
	got, err := io.ReadAll(res.Body)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.JSONEq(want, string(got), "Incorrect response")
}

func (s *testSuite) TestEchoHandlerErrorResponse() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "localhost:8080/echo", nil)
	s.handlers.echoHandler(w, req)

	res := w.Result()
	s.Equal(http.StatusBadRequest, res.StatusCode)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
