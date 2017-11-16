package spark

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RESTTestSuite struct {
	suite.Suite
}

func TestRESTTestSuite(t *testing.T) {
	suite.Run(t, new(RESTTestSuite))
}

type HTTPClientMock struct {
	mock.Mock
}

func (m *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (suite *RESTTestSuite) TestGet() {
	test := func(params map[string]string, expectedURL string) {
		suite.T().Helper()

		httpClient := &HTTPClientMock{}
		expectedRequest, _ := http.NewRequest(
			"GET",
			expectedURL,
			nil,
		)
		expectedRequest.Header.Set("Authorization", "Bearer soopersecret")
		expectedRequest.Header.Set("Content-Type", "application/json")
		httpClient.On(
			"Do",
			expectedRequest,
		).Return(
			&http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(
					bytes.NewBuffer([]byte(
						`{"foo": "bar"}`,
					)),
				),
			},
			nil,
		)
		rest := restClient{
			http:  httpClient,
			url:   "http://example.com/",
			token: "soopersecret",
		}

		var out struct {
			Foo string
		}
		err := rest.Get("resource/path", params, &out)
		suite.Nil(err)
		suite.Equal(out.Foo, "bar")
	}

	test(map[string]string{},
		"http://example.com/resource/path")
	test(map[string]string{"foo": "bar"},
		"http://example.com/resource/path?foo=bar")
	// TODO: test in a way that is independent of map order
	test(map[string]string{"foo": "bar", "biz": "baz"},
		"http://example.com/resource/path?biz=baz&foo=bar")
}

func (suite *RESTTestSuite) TestDelete() {
	test := func(params map[string]string, expectedURL string) {
		suite.T().Helper()

		httpClient := &HTTPClientMock{}
		expectedRequest, _ := http.NewRequest(
			"DELETE",
			expectedURL,
			nil,
		)
		expectedRequest.Header.Set("Authorization", "Bearer soopersecret")
		expectedRequest.Header.Set("Content-Type", "application/json")
		httpClient.On(
			"Do",
			expectedRequest,
		).Return(
			&http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(
					bytes.NewBuffer([]byte(
						`{"foo": "bar"}`,
					)),
				),
			},
			nil,
		)
		rest := restClient{
			http:  httpClient,
			url:   "http://example.com/",
			token: "soopersecret",
		}

		var out struct {
			Foo string
		}
		err := rest.Delete("resource/path", params, &out)
		suite.Nil(err)
		suite.Equal(out.Foo, "bar")
	}

	test(map[string]string{},
		"http://example.com/resource/path")
	test(map[string]string{"foo": "bar"},
		"http://example.com/resource/path?foo=bar")
	// TODO: test in a way that is independent of map order
	test(map[string]string{"foo": "bar", "biz": "baz"},
		"http://example.com/resource/path?biz=baz&foo=bar")
}

func (suite *RESTTestSuite) TestPost() {
	httpClient := &HTTPClientMock{}
	expectedRequest, _ := http.NewRequest(
		"POST",
		"http://example.com/resource/path",
		bytes.NewReader([]byte(`{"foo":"bar"}`)),
	)
	expectedRequest.Header.Set("Authorization", "Bearer soopersecret")
	expectedRequest.Header.Set("Content-Type", "application/json")
	httpClient.On(
		"Do",
		mock.MatchedBy(func(req *http.Request) bool {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return false
			}
			return req.Method == "POST" &&
				req.URL.String() == "http://example.com/resource/path" &&
				string(body) == `{"Foo":"bar"}`
		}),
	).Return(
		&http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(
				bytes.NewBuffer([]byte(
					`{"foo": "bar"}`,
				)),
			),
		},
		nil,
	)
	rest := restClient{
		http:  httpClient,
		url:   "http://example.com/",
		token: "soopersecret",
	}

	in := struct{ Foo string }{Foo: "bar"}

	var out struct {
		Foo string
	}
	err := rest.Post("resource/path", in, &out)
	suite.Nil(err)
	suite.Equal(out.Foo, "bar")
}
