package handlers

import (
	"fmt"
	"strings"
)

func responseStatusLineGenerator(statusCode string) string {
	return fmt.Sprintf("HTTP/1.1 %s\r\n", statusCode)
}

func responseHeaderGenerator(headers map[string]string) string {
	var builder strings.Builder
	for key, value := range headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	builder.WriteString("\r\n")
	return builder.String()
}