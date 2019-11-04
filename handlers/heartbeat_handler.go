package handlers

import (
	"encoding/base64"
	"net/http"
	"strings"
)

var (
	heartbeatGifResponse, _ = base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=")
	heartbeatTextResponse   = []byte("OK")
)

type HeartbeatHandler struct{}

func NewHeartbeatHandler() *HeartbeatHandler {
	handler := &HeartbeatHandler{}
	return handler
}

func (handler *HeartbeatHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch {
	case strings.HasSuffix(req.URL.Path, ".gif"):
		writer.Write(heartbeatGifResponse)
	default:
		writer.Write(heartbeatTextResponse)
	}
}
