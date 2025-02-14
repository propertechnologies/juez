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

	multipartRequest := juez.NewMultiPartRequest[any](e)
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

	multipartRequest := juez.NewMultiPartRequest[any](e)
	multipartRequest.URL("/upload").
		AddFile("file", "file.pdf", reader).
		WithHeaders(map[string]string{"Authorization": "Bearer token"}).
		POST(nil).
		Expect(200)
}

type MockFile struct {
	content string
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	copy(p, m.content)
	return len(m.content), io.EOF
}
