package main

import (
	"fmt"
	"strings"
	"compress/gzip"
)

type Compressor interface {
	Compress(response string) (string, error)
}

type GzipCompressor struct{}

func (c *GzipCompressor) Compress(response string) (string, error) {
	var compressedResponse strings.Builder
	writer := gzip.NewWriter(&compressedResponse)

	_, err := writer.Write([]byte(response))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	return compressedResponse.String(), nil
}

var supportedSchemes = map[string]Compressor{
	"gzip": &GzipCompressor{},
}

func GetCompressor(scheme string) (Compressor, error) {
	compressor, exists := supportedSchemes[scheme]
	if !exists {
		return nil, fmt.Errorf("unsupported compression scheme: %s", scheme)
	}
	return compressor, nil
}