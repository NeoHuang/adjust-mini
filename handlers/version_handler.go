package handlers

import (
	"fmt"
	"net/http"
)

type VersionHandler struct {
	version string
}

func NewVersionHandler(version string) *VersionHandler {
	handler := &VersionHandler{
		version: version,
	}
	return handler
}

func (handler *VersionHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"version": "%s"}`, handler.version)
}
