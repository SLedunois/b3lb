package admin

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/utils"
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddInstance(t *testing.T) {
	type test struct {
		name           string
		body           io.Reader
		expectedStatus int
		expectedBody   string
	}

	tests := []test{
		{
			name:           "Nil body should trigger a bad request status",
			body:           nil,
			expectedStatus: 400,
			expectedBody:   "",
		},
		{
			name:           "Correct body should create a bbb instance",
			body:           bytes.NewBufferString(`{"url": "http://localhost/bigbluebutton", "secret": "supersecret"}`),
			expectedStatus: 201,
			expectedBody:   `{"url":"http://localhost/bigbluebutton","secret":"supersecret"}`,
		},
		{
			name:           "Adding a duplication should returns a 409 conflict status",
			body:           bytes.NewBufferString(`{"url": "http://localhost/bigbluebutton", "secret": "supersecret"}`),
			expectedStatus: 409,
			expectedBody:   "",
		},
	}

	headers := map[string]string{
		"Authorization": TestUtil.DefaultAPIKey(),
		"Content-Type":  "application/json",
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := TestUtil.ExecuteRequestWithHeaders(router, "POST", "/admin/servers", test.body, headers)
			assert.Equal(t, test.expectedStatus, w.Code)

			if test.expectedBody != "" {
				assert.Equal(t, test.expectedBody, w.Body.String())
			}
		})
	}
}

func TestListInstances(t *testing.T) {
	headers := map[string]string{
		"Authorization": TestUtil.DefaultAPIKey(),
	}

	w := TestUtil.ExecuteRequestWithHeaders(router, "GET", "/admin/servers", nil, headers)
	assert.Equal(t, 200, w.Code)
	var arr []string
	err := json.Unmarshal(w.Body.Bytes(), &arr)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, utils.ArrayContainsString(arr, "http://localhost/bigbluebutton"), true)
}

func TestDeleteInstance(t *testing.T) {

	type test struct {
		name           string
		url            string
		expectedStatus int
	}

	tests := []test{
		{
			name:           "Delete without an url should return a 400",
			url:            "/admin/servers",
			expectedStatus: 400,
		},
		{
			name:           "Delete a non existing instance should return a 404",
			url:            "/admin/servers?url=http://fakebbb",
			expectedStatus: 404,
		},
		{
			name:           "Delete an existing instance should return a 204",
			url:            "/admin/servers?url=http://localhost/bigbluebutton",
			expectedStatus: 204,
		},
	}

	headers := map[string]string{
		"Authorization": TestUtil.DefaultAPIKey(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := TestUtil.ExecuteRequestWithHeaders(router, "DELETE", test.url, nil, headers)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}