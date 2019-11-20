package handlers

import (
	"fmt"
	"net/http"
	"os"
)

type VersionHandler struct {
	version string
	podName string
}

func NewVersionHandler(version string) *VersionHandler {
	podName := os.Getenv("POD_NAME")
	handler := &VersionHandler{
		version: version,
		podName: podName,
	}
	return handler
}

func (handler *VersionHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"version": "%s", "pod_name":"%s"}`, handler.version, handler.podName)
}
