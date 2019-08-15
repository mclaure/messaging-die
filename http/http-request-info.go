package http

import (
	Json "encoding/json"
	"fmt"
	NetHttp "net/http"
	"strings"
)

// RequestInfo ...
type RequestInfo struct {
	Host          string
	APIPattern    string
	Method        string
	Body          []byte
	URLParameters map[string]string
}

func createParameters(requestURI string) map[string]string {
	var m map[string]string

	paramsT := strings.Split(requestURI, "?")
	if len(paramsT) == 2 {
		m = make(map[string]string)
		params := strings.Split(paramsT[1], "&")
		for _, param := range params {
			pair := strings.Split(param, "=")
			if len(pair) == 2 {
				m[pair[0]] = pair[1]
			}
		}
	}

	return m
}

func createRequestInfo(request *NetHttp.Request) RequestInfo {
	return RequestInfo{
		Host:          request.Host,
		APIPattern:    request.URL.Path,
		Method:        request.Method,
		Body:          nil, // TODO. For POST Methods
		URLParameters: createParameters(request.RequestURI),
	}
}

// GetURLParameters ...
func GetURLParameters(params map[string]string) string {
	var urlParameters string
	pairs := []string{}
	for key, value := range params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}

	if len(pairs) > 0 {
		urlParameters = strings.Join(pairs, "&")
	}

	return urlParameters
}

// GetRequestInfoBytes ...
func GetRequestInfoBytes(request *NetHttp.Request) []byte {
	requestInfo := createRequestInfo(request)
	b, _ := Json.Marshal(requestInfo)

	return b
}

// GetRequestInfo ...
func GetRequestInfo(bytes []byte) RequestInfo {
	var requestInfo RequestInfo
	Json.Unmarshal(bytes, &requestInfo)

	return requestInfo
}
