package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	var responseTests = []struct {
		name string
		statusIn int
		dataIn interface{}
		errIn error
		expectedStatus int
		expectedError string
		expectedData interface{}
	}{
		{
			"error of type ResponseError",
			http.StatusBadRequest,
			nil,
			&testError {
					httpStatus: http.StatusNotFound,
					kind:       "InvalidUserInput",
					err:        errors.New("validation stuff failed"),
				},
			http.StatusNotFound,
			"InvalidUserInput: validation stuff failed",
			nil,
		},
		{
			"error of type error with data",
			http.StatusBadRequest,
			map[string]interface{}{"test":"test"},
			errors.New("err"),
			http.StatusInternalServerError,
			"err",
			map[string]interface{}{"test":"test"},
		},
		{
			"no error with data",
			http.StatusOK,
			map[string]interface{}{"test":"test"},
			nil,
			http.StatusOK,
			"",
			map[string]interface{}{"test":"test"},
		},
	}
	for _, tt := range responseTests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			// error present and of type response error
			RespondWithJSON(rec, tt.statusIn, tt.dataIn, tt.errIn)
			result := rec.Result()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))
			var response ResponseBody
			b, _ := ioutil.ReadAll(result.Body)
			_ = json.Unmarshal(b, &response)
			if tt.expectedError == "" {
				assert.Nil(t, response.ErrorMessage)
			} else {
				assert.Equal(t, tt.expectedError, *response.ErrorMessage)
			}
			assert.Equal(t, tt.expectedData, response.Data)
		})
	}
}

func TestRespondWithJSONError(t *testing.T) {
	rec := httptest.NewRecorder()
	RespondWithJSONError(rec, &testError{
		httpStatus: http.StatusForbidden,
		kind:       "KUndefined",
		err:        errors.New("error"),
	})
	assert.Equal(t, rec.Result().StatusCode, http.StatusForbidden)
	body := readBody(rec.Result().Body)
	assert.Equal(t, *body.ErrorMessage, "KUndefined: error")
	assert.Nil(t, body.Data)
}

func TestRespondWithJSONData(t *testing.T) {
	rec := httptest.NewRecorder()
	RespondWithJSONData(rec, http.StatusOK, map[string]string{"test":"data"})
	assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
	body := readBody(rec.Result().Body)
	assert.Nil(t, body.ErrorMessage)
	data := body.Data.(map[string]interface{})
	assert.Equal(t, data["test"], "data")
}

func readBody(body io.ReadCloser) ResponseBody {
	var rb ResponseBody
	b, _ := ioutil.ReadAll(body)
	_ = json.Unmarshal(b, &rb)
	return rb
}

type testError struct {
	httpStatus int
	kind string
	err error
}

func (te *testError) HTTPStatusCode() int {
	return te.httpStatus
}

func (te *testError) ErrorKind() string {
	return te.kind
}

func (te *testError) Error() string {
	return fmt.Sprintf("%s: %s", te.kind, te.err.Error())
}
