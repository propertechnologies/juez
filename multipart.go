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
	multipartRequest[T any, R any] struct {
		baseRequestResponse
		writer *multipart.Writer
		buffer *bytes.Buffer
	}
)

func NewMultiPartRequest[T any](e *gin.Engine) *multipartRequest[T, T] {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	return &multipartRequest[T, T]{
		baseRequestResponse: baseRequestResponse{
			engine: e,
		},
		writer: writer,
		buffer: &buffer,
	}
}

func (r *multipartRequest[T, R]) POST(b T) *multipartRequest[T, R] {
	if err := r.writer.Close(); err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, r.url, r.buffer)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", r.writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	r.engine.ServeHTTP(recorder, req)

	r.responseRecorder = recorder

	r.writer.Close()

	return r
}

func (r *multipartRequest[T, R]) AddFormData(b T, name string, value string) *multipartRequest[T, R] {
	if err := r.writer.WriteField(name, value); err != nil {
		panic(err)
	}

	return r
}

func (r *multipartRequest[T, R]) AddFile(fieldName string, fileName string, reader io.Reader) *multipartRequest[T, R] {
	// Add file
	fileWriter, err := r.writer.CreateFormFile(fieldName, fieldName)
	if err != nil {
		panic(err)
	}

	// Copy the file content to the file writer
	if _, err = io.Copy(fileWriter, reader); err != nil {
		panic(err)
	}

	return r
}

func (r *multipartRequest[T, R]) URL(u string) *multipartRequest[T, R] {
	r.url = u

	return r
}
