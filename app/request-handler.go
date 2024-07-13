package main

func HandleRequest(request Request) string {
	return requestHandlerFactory(request.RequestLine.Method)(request)
}

type requestHandler func (request Request) string

func requestHandlerFactory(methodType string) requestHandler {
	handlerMap := map[string]requestHandler{
		"GET": handleGetRequest,
		"POST": handlePostRequest,
	}

	if handler, ok := handlerMap[methodType]; ok {
		return handler
	}

	return nil // TODO: handle other methods
}