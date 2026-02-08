package parser

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Request struct {
	Method string
	Path   string
	Body   string
}

const (
	GET  = "GET"
	POST = "POST"
)

var allowedMethods = []string{GET, POST}

const reqFirstLnSep = " "
const minReqItems = 2 // Only the method and path will be extracted from the request
const pathPrefix = "/"
const reqBodySep = "\r\n\r\n"
const contentLengTxt = "content-length"

func setMethod(method string, req *Request) error {
	if !slices.Contains(allowedMethods, method) {
		return fmt.Errorf("invalid method: %s (allowed: %v)", method, allowedMethods)
	}
	req.Method = method
	return nil
}

func setPath(path string, req *Request) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid path: %s (must start with %s)", path, pathPrefix)
	}
	req.Path = path
	return nil
}

func getContentLength(rawReq string) (int, error) {
	lines := strings.Split(rawReq, "\n")
	for _, line := range lines {
		lowerLine := strings.ToLower(strings.TrimSpace(line))
		if !strings.Contains(lowerLine, contentLengTxt) {
			continue
		}
		lineParts := strings.Split(lowerLine, ":")
		if len(lineParts) < 2 {
			return -1, fmt.Errorf("malformed request: %s (must contain a valid content length)", rawReq)
		}
		contentLen, err := strconv.Atoi(strings.TrimSpace(lineParts[1]))
		if err != nil {
			return -1, err
		}
		return contentLen, nil
	}
	return -1, fmt.Errorf("malformed request: %s (must contain a content length for the body)", rawReq)
}

func SetRequestData(rawReq string, req *Request) (*Request, error) {
	reqFirstLineTokens := strings.SplitN(rawReq, reqFirstLnSep, 3)
	if len(reqFirstLineTokens) < minReqItems {
		return req, fmt.Errorf("malformed request: %s (must contain a verb and a path)", rawReq)
	}
	if err := setMethod(reqFirstLineTokens[0], req); err != nil {
		return req, err
	}
	if err := setPath(strings.TrimRight(reqFirstLineTokens[1], " \r\n"), req); err != nil {
		return req, err
	}
	if req.Method == POST {
		reqHeadersBody := strings.SplitN(rawReq, reqBodySep, 2)
		if len(reqHeadersBody) == 1 { // If the body is empty
			return req, nil
		}
		contentLen, err := getContentLength(rawReq)
		if err != nil {
			return nil, err
		}
		req.Body = reqHeadersBody[1][:contentLen]
	}
	return req, nil
}
