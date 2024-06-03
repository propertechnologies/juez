package juez

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// requestResponse is a helper struct that represents an HTTP request and response.
type requestResponse[T any, R any] struct {
	u                string
	responseRecorder *httptest.ResponseRecorder
	engine           *gin.Engine
	headers          map[string]string
}

// URL sets the URL for the request.
func (r *requestResponse[T, R]) URL(u string) *requestResponse[T, R] {
	r.u = u

	return r
}

// GET sends a GET request to the specified URL and records the response.
func (r *requestResponse[T, R]) GET() *requestResponse[T, R] {
	req := r.newRequest(http.MethodGet, r.u, nil)

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

// NewRequestWithResponse creates a new requestResponse instance with the specified gin.Engine.
func NewRequestWithResponse[T any, R any](e *gin.Engine) *requestResponse[T, R] {
	return &requestResponse[T, R]{
		engine: e,
	}
}

// newRequest creates a new http.Request with the specified method, URL, and body.
func (r *requestResponse[T, R]) newRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}

	for header, value := range r.headers {
		req.Header.Set(header, value)
	}

	return req
}

// Expect checks if the HTTP response status code matches the expected status code.
// If the status codes do not match, it panics with an error message.
func (r *requestResponse[T, R]) Expect(httpStatus int) *requestResponse[T, R] {
	if httpStatus != r.responseRecorder.Code {
		panic(
			fmt.Sprintf(
				"expected:%d and received: %d. %s\n",
				httpStatus,
				r.responseRecorder.Code,
				r.responseRecorder.Body,
			),
		)
	}

	return r
}

// Body returns the response body as the specified type.
func (r *requestResponse[T, R]) Body() R {
	b := r.responseRecorder.Body.Bytes()
	return BodyToReceive[R](b)
}

// BodyToReceive converts the response body bytes to the specified type using JSON unmarshaling.
func BodyToReceive[T any](b []byte) T {
	var actual T

	if len(b) == 0 {
		return actual
	}

	if err := json.Unmarshal(b, &actual); err != nil {
		panic(err)
	}

	return actual
}

func (r *requestResponse[T, R]) POST(b T) *requestResponse[T, R] {
	body := bodyToSend(&b)
	req := r.newRequest(http.MethodPost, r.u, bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func bodyToSend[T any](d *T) []byte {
	r, err := json.Marshal(d)
	if err != nil {
		fmt.Printf("error: %v", err)
		return []byte{}
	}

	return r
}
