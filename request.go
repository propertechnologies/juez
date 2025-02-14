package juez

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type (
	requestResponse[T any, R any] struct {
		baseRequestResponse
	}

	db interface {
		Rollback() error
	}

	HttpServer interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
	}
)

func RunIntegrationTest(tx db, f func()) {
	defer func() {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
	}()

	f()
}

func NewRequest[T any](e HttpServer) *requestResponse[T, T] {
	return &requestResponse[T, T]{
		baseRequestResponse: baseRequestResponse{
			engine: e,
		},
	}
}

func NewRequestWithResponse[T any, R any](e HttpServer) *requestResponse[T, R] {
	return &requestResponse[T, R]{
		baseRequestResponse: baseRequestResponse{
			engine: e,
		},
	}
}

func (r *requestResponse[T, R]) URL(u string) *requestResponse[T, R] {
	r.url = u

	return r
}

func (r *requestResponse[T, R]) POST(b T) *requestResponse[T, R] {
	body := bodyToSend(&b)
	req := r.newRequest(http.MethodPost, r.url, bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) POSTWithJson(body []byte) *requestResponse[T, R] {
	req := r.newRequest(http.MethodPost, r.url, bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) PUT(b T) *requestResponse[T, R] {
	body := bodyToSend(&b)
	req := r.newRequest(http.MethodPut, r.url, bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) GET() *requestResponse[T, R] {
	req := r.newRequest(http.MethodGet, r.url, nil)

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) WithHeaders(headers map[string]string) *requestResponse[T, R] {
	r.headers = headers
	return r
}

func (r *requestResponse[T, R]) DELETE() *requestResponse[T, R] {
	req := r.newRequest(http.MethodDelete, r.url, nil)

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) PATCH(b T) *requestResponse[T, R] {
	body := bodyToSend(&b)
	req := r.newRequest(http.MethodPatch, r.url, bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) PATCHWithJson(body []byte) *requestResponse[T, R] {
	req := r.newRequest(http.MethodPatch, r.url, bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	return r
}

func (r *requestResponse[T, R]) Body() R {
	b := r.responseRecorder.Body.Bytes()
	return BodyToReceive[R](b)
}

func (r *requestResponse[T, R]) BodyBytes() []byte {
	return r.responseRecorder.Body.Bytes()
}

func (r *requestResponse[T, R]) Expect(httpStatus int) *requestResponse[T, R] {
	r.baseRequestResponse.Expect(httpStatus)

	return r
}

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

func bodyToSend[T any](d *T) []byte {
	r, err := json.Marshal(d)
	if err != nil {
		fmt.Printf("error: %v", err)
		return []byte{}
	}

	return r
}
