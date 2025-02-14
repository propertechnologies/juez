package juez

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type (
	multipartRequestResponse[T any, R any] struct {
		baseRequestResponse
		writer *multipart.Writer
		buffer *bytes.Buffer
	}
)

func NewMultiPartRequestResponse[T any, R any](e *gin.Engine) *multipartRequestResponse[T, R] {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	return &multipartRequestResponse[T, R]{
		baseRequestResponse: baseRequestResponse{
			engine: e,
		},
		writer: writer,
		buffer: &buffer,
	}
}

func (r *multipartRequestResponse[T, R]) POST(b T) *multipartRequestResponse[T, R] {
	req, err := http.NewRequest(http.MethodPost, r.url, r.buffer)
	if err != nil {
		panic(err)
	}

	for header, value := range r.headers {
		req.Header.Set(header, value)
	}

	req.Header.Set("Content-Type", r.writer.FormDataContentType())

	if err := r.writer.Close(); err != nil {
		panic(err)
	}

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	r.writer.Close()

	return r
}

func (r *multipartRequestResponse[T, R]) AddFormData(b T, name string, value string) *multipartRequestResponse[T, R] {
	if err := r.writer.WriteField(name, value); err != nil {
		panic(err)
	}

	return r
}

func (r *multipartRequestResponse[T, R]) AddFile(fieldName string, fileName string, reader io.Reader) *multipartRequestResponse[T, R] {
	fileWriter, err := r.writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		panic(err)
	}

	if _, err = io.Copy(fileWriter, reader); err != nil {
		panic(err)
	}

	return r
}

func (r *multipartRequestResponse[T, R]) URL(u string) *multipartRequestResponse[T, R] {
	r.url = u

	return r
}

func (r *multipartRequestResponse[T, R]) Body() R {
	b := r.responseRecorder.Body.Bytes()
	return BodyToReceive[R](b)
}

func (r *multipartRequestResponse[T, R]) Expect(httpStatus int) *multipartRequestResponse[T, R] {
	r.baseRequestResponse.Expect(httpStatus)

	return r
}

func (r *multipartRequestResponse[T, R]) WithHeaders(headers map[string]string) *multipartRequestResponse[T, R] {
	r.baseRequestResponse.WithHeaders(headers)

	return r
}
