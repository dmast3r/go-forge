package main

import "strings"

type RequestLine struct {
	Method string
	Path   string
	HttpVersion string
}

type Request struct {
	RequestLine RequestLine
	Headers map[string]string
	Body string
}

func ParseRequest(request string) Request {
	requestLine, rest, _ := strings.Cut(request, "\r\n")
	headers, body, _ := strings.Cut(rest, "\r\n\r\n")

	return Request{
		RequestLine: parseRequestLine(requestLine),
		Headers: parseHeaders(headers),
		Body: body,
	}
}

func parseRequestLine(requestLine string) RequestLine {
	method, rest, _ := strings.Cut(requestLine, " ")
	path, httpVersion, _ := strings.Cut(rest, " ")

	return RequestLine{
		Method: method,
		Path: path,
		HttpVersion: httpVersion,
	}
}

func parseHeaders(headers string) map[string]string {
	headerLines := strings.Split(headers, "\r\n")
	headersMap := make(map[string]string)

	for _, headerLine := range headerLines {
		headerName, headerValue, _ := strings.Cut(headerLine, ": ")
		headersMap[headerName] = headerValue
	}
	
	return headersMap
}