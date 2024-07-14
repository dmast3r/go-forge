package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmast3r/go-forge/app/utils"
)

func handleGetRequest(request Request) string {
	pathString := request.RequestLine.Path
	paths := strings.Split(pathString, "/")
	compressionSchemes := request.Headers["Accept-Encoding"]

	if pathString == "/" || pathString == "/ping" {
		return pingResponseGenerator(compressionSchemes)
	} else if paths[1] == "echo" && len(paths) > 2 {
		return echoResponseGenerator(paths[2], compressionSchemes)
	} else if paths[1] == "user-agent" {
		return userAgentResponseGenerator(request.Headers["User-Agent"], false, compressionSchemes)
	} else if paths[1] == "files" && len(paths) > 2 {
		return fileResponseGenerator(paths[2], false, compressionSchemes)
	}

	return "HTTP/1.1 404 Not Found\r\n\r\n" // TODO: handle other paths
}

func pingResponseGenerator(compressionSchemes string) string {
	responseHeaders := map[string]string{}

	return handleResponse("200 OK", "", compressionSchemes, responseHeaders)
}

func echoResponseGenerator(echoString string, compressionSchemes string) string {
	responseHeaders := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(echoString)),
	}

	return handleResponse("200 OK", echoString, compressionSchemes, responseHeaders)
}

func userAgentResponseGenerator(userAgent string, needsCompression bool, compressionSchemes string) string {
	responseHeaders := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(userAgent)),
	}

	return handleResponse("200 OK", userAgent, compressionSchemes, responseHeaders)
}

func fileResponseGenerator(fileName string, needsCompression bool, compressionSchemes string) string {
	filePath := filepath.Join(os.Getenv("WORKING_DIRECTORY"), fileName)
	content, err := os.ReadFile(filePath)

	if err != nil {
		return handleResponse("404 Not Found", "", compressionSchemes, map[string]string{})
	}

	contentString := string(content)
	responseHeaders := map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": fmt.Sprintf("%d", len(contentString)),
	}

	return handleResponse("200 OK", contentString, compressionSchemes, responseHeaders)
}

func handleResponse(statusCode, response, compressionSchemes string, responseHeaders map[string]string) string {
	compressionScheme, compressor := getCompressor(compressionSchemes)
	if compressor == nil {
		return responseStatusLineGenerator(statusCode) + responseHeaderGenerator(responseHeaders) + response
	}

	compressedResponse, err := compressor.Compress(response)
	if err != nil {
		return responseStatusLineGenerator(statusCode) + responseHeaderGenerator(responseHeaders) + response
	}

	responseHeaders["Content-Encoding"] = compressionScheme
	responseHeaders["Content-Length"] = fmt.Sprintf("%d", len(compressedResponse))
	return responseStatusLineGenerator(statusCode) + responseHeaderGenerator(responseHeaders) + compressedResponse
}

func getCompressor(compressionSchemes string) (string, utils.Compressor) {
	for _, compressionScheme := range strings.Split(compressionSchemes, ", ") {
		if compressor, err := utils.GetCompressor(compressionScheme); err == nil {
			return compressionScheme, compressor
		}
	}

	return "", nil
}
