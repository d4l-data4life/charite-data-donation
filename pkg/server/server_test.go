package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/d4l-data4life/charite-data-donation/internal/testutils"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
)

func TestCors(t *testing.T) {
	tests := []struct {
		name                  string
		reply                 *httptest.ResponseRecorder
		request               *http.Request
		requestHeader         string // one header to include in request (cannot use maps here)
		requestHeaderContent  string // header value
		expectHeaders         bool   // whether expectedHeader should be present in reply
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:                  "Access-Control-Allow-Origin header should be present",
			reply:                 httptest.NewRecorder(),
			request:               httptest.NewRequest("GET", "/checks/liveness", nil),
			requestHeader:         "Origin",
			requestHeaderContent:  "localhost",
			expectHeaders:         true,
			expectedHeader:        "Access-Control-Allow-Origin",
			expectedHeaderContent: "localhost",
		},
		{
			name:                  "Access-Control-Expose-Headers header should be present",
			reply:                 httptest.NewRecorder(),
			request:               httptest.NewRequest("GET", "/checks/liveness", nil),
			requestHeader:         "Origin",
			requestHeaderContent:  "localhost",
			expectHeaders:         true,
			expectedHeader:        "Access-Control-Expose-Headers",
			expectedHeaderContent: "Link, X-Csrf-Token",
		},
		{
			name:                  "Access-Control-Allow-Credentials header should be present",
			reply:                 httptest.NewRecorder(),
			request:               httptest.NewRequest("GET", "/checks/liveness", nil),
			requestHeader:         "Origin",
			requestHeaderContent:  "localhost",
			expectHeaders:         true,
			expectedHeader:        "Access-Control-Allow-Credentials",
			expectedHeaderContent: "true",
		},
		{
			name:                  "Origin matches not",
			reply:                 httptest.NewRecorder(),
			request:               httptest.NewRequest("GET", "/checks/liveness", nil),
			requestHeader:         "Origin",
			requestHeaderContent:  "http://www.data4life.care",
			expectHeaders:         false,
			expectedHeader:        "Access-Control-Allow-Origin",
			expectedHeaderContent: "localhost",
		},
	}

	server := testutils.GetTestMockServer()
	defer models.GetDB().Close()

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			test.request.Header.Set(test.requestHeader, test.requestHeaderContent)

			server.Mux().ServeHTTP(test.reply, test.request)
			if test.expectHeaders {
				assert.Equal(t, test.expectedHeaderContent, test.reply.Header().Get(test.expectedHeader))
			} else {
				assert.Equal(t, "", test.reply.Header().Get(test.expectedHeader))
			}
		})
	}
}
