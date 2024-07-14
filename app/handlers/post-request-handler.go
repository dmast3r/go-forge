package handlers

import (
	"os"
	"path/filepath"
	"strings"
)

func handlePostRequest(request Request) string {
	paths := strings.Split(request.RequestLine.Path, "/")

	if paths[1] == "files" && len(paths) > 2 {
		return handleFileUploadRequest(request)
	}
	return "HTTP/1.1 404 Not Found\r\n\r\n"
}

func handleFileUploadRequest(request Request) string {
	fileName := strings.Split(request.RequestLine.Path, "/")[2]
	filePath := filepath.Join(os.Getenv("WORKING_DIRECTORY"), fileName)
	content := request.Body

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return responseStatusLineGenerator("500 Internal Server Error") + responseHeaderGenerator(map[string]string{})
	}

	return responseStatusLineGenerator("201 Created") + responseHeaderGenerator(map[string]string{})
}