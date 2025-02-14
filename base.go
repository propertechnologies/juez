package juez

import (
	"fmt"
	"net/http/httptest"
)

type (
	baseRequestResponse struct {
		url              string
		responseRecorder *httptest.ResponseRecorder
		engine           HttpServer
		headers          map[string]string
	}
)

func (r *baseRequestResponse) Expect(httpStatus int) *baseRequestResponse {
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
