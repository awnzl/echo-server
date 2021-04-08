package endpointhandlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/awnzl/echo-server/internal/logger"
)

type testSuite struct{
	suite.Suite
	log *zap.Logger
}

func (s *testSuite) SetupTest() {
	s.log = logger.NewZap("info")
}

func (s *testSuite) TearDownSuite() {
	s.log.Sync()
}

func (s *testSuite) TestNameHandler() {
	want := nameResponse{
		ServiceName: "echo-server",
	}
	handlers := New(s.log)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "localhost:8080", nil)
	handlers.nameHandler(w, req)

	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(s.T(), err)

	got := nameResponse{}
	if err := json.Unmarshal(data, &got); err != nil {
		s.log.Error(err.Error())
	}
	assert.Equal(s.T(), want, got, "Incorrect response")

}

func (s *testSuite) TestEchoHandler() {
	want := echoResponse{
		EchoWord: "test echo word",
	}
	handlers := New(s.log)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST",
		"localhost:8080/echo",
		strings.NewReader(`{"word":"test echo word"}`),
	)
	handlers.echoHandler(w, req)

	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(s.T(), err)

	got := echoResponse{}
	if err := json.Unmarshal(data, &got); err != nil {
		s.log.Error(err.Error())
	}
	assert.Equal(s.T(), want, got, "Incorrect response")
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
