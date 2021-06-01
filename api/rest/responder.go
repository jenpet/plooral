package rest

import (
	"fmt"
	"net/http"
)

const ContentTypeUTF8JSON = "application/json; charset=utf-8"

// RespondWithJSONData convenience wrapper for RespondWithJSON
func RespondWithJSONData(w http.ResponseWriter, httpStatus int, data interface{}) {
	RespondWithJSON(w, httpStatus, data, nil)
}

// RespondWithJSONError convenience wrapper for RespondWithJSON
func RespondWithJSONError(w http.ResponseWriter, err error) {
	RespondWithJSON(w, 0, nil, err)
}

// RespondWithJSON write the given data into the body and set the HTTP status if no error is provided.
// In case an error is present its message will be used as an error message in the response and the response HTTP status
// will be overwritten with HTTP 500.
// If the present error can be asserted into a ResponseError it will use the different arguments to build up the error message
// and extract the underlying status if it is more specific than HTTP 500.
func RespondWithJSON(w http.ResponseWriter, httpStatus int, data interface{}, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		// if an error occurs the default http status is http 500 (internal server error)
		// In case error is of type ResponseError it might override the existing error code
		httpStatus = http.StatusInternalServerError
	}
	// assert a ResponseError
	respErr,_ := err.(ResponseError)
	// if type ResponseError is present check whether it has a more critical (i.e. higher) HTTP status code
	if respErr != nil {
		if respErr.HTTPStatusCode() < httpStatus {
			httpStatus = respErr.HTTPStatusCode()
		}
		errorMsg = responseErrorToString(respErr)
	}

	w.Header().Add("Content-Type", ContentTypeUTF8JSON)
	w.WriteHeader(httpStatus)
	resp := ResponseBody{ Data: data }
	if errorMsg != "" {
		resp.SetErrorMessage(errorMsg)
	}
	_,_ = w.Write(resp.JSON())
}

func responseErrorToString(rErr ResponseError) string {
	return fmt.Sprintf("%s: %s", rErr.ErrorKind(), rErr.Error())
}

// ResponseError provides information required to return a more detailed HTTP response.
type ResponseError interface {
	// HTTPStatusCode returns the HTTP status code which should be used
	HTTPStatusCode() int

	// ErrorKind returns a simple classification of the error
	ErrorKind() string

	// Error returns the detailed error message
	Error() string
}