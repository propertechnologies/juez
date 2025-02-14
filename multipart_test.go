package juez_test

import (
	"io"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/propertechnologies/juez"
)

func TestMultipart(t *testing.T) {
	e := gin.Default()
	e.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if file == nil || err != nil {
			c.String(400, "file not found or error %s", err)
			return
		}
		if file.Size == 0 {
			c.String(400, "file is empty")
			return
		}
		c.String(200, file.Filename)
	})

	reader := &MockFile{
		content: "hello world, this is a file",
	}

	multipartRequest := juez.NewMultiPartRequestResponse[any, any](e)
	multipartRequest.URL("/upload").AddFile("file", "file.txt", reader).POST(nil).Expect(200)
}

func TestMultipartWithHeaders(t *testing.T) {
	e := gin.Default()
	e.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if file == nil || err != nil {
			c.String(400, "file not found or error %s", err)
			return
		}

		if !strings.HasSuffix(file.Filename, ".pdf") {
			c.String(400, "file is not a pdf")
			return
		}

		if file.Size == 0 {
			c.String(400, "file is empty")
			return
		}

		if c.GetHeader("Authorization") != "Bearer token" {
			c.String(401, "token is missing")
			return
		}

		c.String(200, file.Filename)
	})

	reader := &MockFile{
		content: "hello world, this is a file",
	}

	multipartRequest := juez.NewMultiPartRequestResponse[any, any](e)
	multipartRequest.URL("/upload").
		AddFile("file", "file.pdf", reader).
		WithHeaders(map[string]string{"Authorization": "Bearer token"}).
		POST(nil).
		Expect(200)
}

func TestMultipartWithResponse(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
	}

	response := Response{
		Message: "hello world",
	}

	e := gin.Default()
	e.POST("/upload", func(c *gin.Context) {

		c.JSON(200, response)
	})

	reader := &MockFile{
		content: "hello world, this is a file",
	}

	multipartRequest := juez.NewMultiPartRequestResponse[any, Response](e)
	resp := multipartRequest.URL("/upload").
		AddFile("file", "file.pdf", reader).
		WithHeaders(map[string]string{"Authorization": "Bearer token"}).
		POST(nil).
		Expect(200).
		Body()

	if resp.Message != "hello world" {
		t.Errorf("expected: hello world and received: %s", resp.Message)
	}
}

type MockFile struct {
	content string
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	copy(p, m.content)
	return len(m.content), io.EOF
}
